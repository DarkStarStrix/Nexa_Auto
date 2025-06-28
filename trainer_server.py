import os
import uuid

import requests
import torch
from datasets import load_dataset
from fastapi import FastAPI, HTTPException, BackgroundTasks
from pydantic import BaseModel
from transformers import AutoModelForCausalLM, AutoTokenizer, TrainingArguments, Trainer, AutoConfig

app = FastAPI(title="Nexa Auto Trainer Backend", version="1.0")

jobs = {}
SESSION_SERVER_URL = "http://127.0.0.1:8765"

def get_token():
    try:
        resp = requests.get(f"{SESSION_SERVER_URL}/get_token")
        if resp.status_code == 200:
            return resp.json()["token"]
    except Exception:
        pass
    return None

class TrainRequest(BaseModel):
    model: str
    dataset: str
    output: str

@app.post("/train")
def start_training(req: TrainRequest, background_tasks: BackgroundTasks):
    job_id = str(uuid.uuid4())
    log_path = f"nexa_output/train_{job_id}.log"
    jobs[job_id] = {"status": "running", "log": log_path}
    background_tasks.add_task(run_training, req.model, req.dataset, req.output, log_path, job_id)
    return {"job_id": job_id}

def run_training(model_name, dataset_name, new_model_name, log_path, job_id):
    os.makedirs("nexa_output", exist_ok=True)
    hf_token = get_token()
    with open(log_path, "w") as logf:
        if not hf_token:
            logf.write("[ERROR] No Hugging Face token found in session. Aborting.\n")
            jobs[job_id]["status"] = "error"
            return
        try:
            logf.write(f"[INFO] Loading model and tokenizer: {model_name}\n")
            logf.flush()
            config = AutoConfig.from_pretrained(model_name, use_auth_token=hf_token)
            model = AutoModelForCausalLM.from_pretrained(model_name, config=config, use_auth_token=hf_token)
            tokenizer = AutoTokenizer.from_pretrained(model_name, use_auth_token=hf_token)
            logf.write(f"[INFO] Loading dataset: {dataset_name}\n")
            logf.flush()
            dataset = load_dataset(dataset_name, split="train")
            def tokenize_function(examples):
                return tokenizer(examples['text'], truncation=True, padding='max_length', max_length=128)
            tokenized_dataset = dataset.map(tokenize_function, batched=True)
            output_dir = os.path.join(os.getcwd(), "nexa_output", new_model_name)
            training_args = TrainingArguments(
                output_dir=output_dir,
                num_train_epochs=1,
                per_device_train_batch_size=2,
                save_steps=10,
                save_total_limit=1,
                logging_steps=5,
                report_to=[],
                push_to_hub=False,
                logging_dir=os.path.join(output_dir, "logs")
            )
            trainer = Trainer(
                model=model,
                args=training_args,
                train_dataset=tokenized_dataset,
                tokenizer=tokenizer,
            )
            logf.write("[INFO] Starting training...\n")
            logf.flush()
            trainer.train()
            trainer.save_model(output_dir)
            tokenizer.save_pretrained(output_dir)
            logf.write(f"[SUCCESS] Model and tokenizer saved to {output_dir}!\n")
            logf.flush()
            del model
            del trainer
            torch.cuda.empty_cache()
            jobs[job_id]["status"] = "finished"
        except Exception as e:
            logf.write(f"[ERROR] {str(e)}\n")
            logf.flush()
            jobs[job_id]["status"] = "error"

@app.get("/logs/{job_id}")
def get_logs(job_id: str):
    job = jobs.get(job_id)
    if not job:
        raise HTTPException(status_code=404, detail="Job not found")
    if not os.path.exists(job["log"]):
        return {"logs": ""}
    with open(job["log"], "r") as f:
        return {"logs": f.read()}

@app.get("/status/{job_id}")
def get_status(job_id: str):
    job = jobs.get(job_id)
    if not job:
        raise HTTPException(status_code=404, detail="Job not found")
    return {"status": job["status"]}

@app.get("/health")
def health():
    return {"status": "ok"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("trainer_server:app", host="0.0.0.0", port=8770, reload=False)
