import numpy as np
import cv2
import joblib
import numpy as np
import seldon_core
from seldon_core.user_model import SeldonComponent
from seldon_core.user_model import SeldonResponse
from typing import Dict, List, Union, Iterable
import os
import logging
import onnx
from urllib.parse import urlparse
from seldon_core.utils import getenv
import onnxruntime as rt
import time
from minio import Minio
import tempfile
import yaml
import random

logger = logging.getLogger(__name__)

class DummyDetector(SeldonComponent):
  def __init__(self, model_uri: str = None,  method: str = "predict", modelUri: str = None, type: str = None):
    
    super().__init__()
    self.model_uri = model_uri
    self.method = method
    self.ready = False
    self.out_dir = tempfile.mkdtemp()
    
    model_file =  os.path.join(self._download_model(self.model_uri, self.out_dir), "model.onnx")

    self._model = model_file
    self.session = rt.InferenceSession(self._model, None)
    self.input_name = self.session.get_inputs()[0].name
    self.output_name = self.session.get_outputs()[0].name
    self.ready = True
    print("init and model loading done!!!")
 
  def init_metadata(self):
    file_path = os.path.join(self._download_model(self.model_uri, self.out_dir), "metadata.yaml")

    try:
      with open(file_path, "r") as f:
        self.mdata = yaml.safe_load(f.read())
        return self.mdata
   
    except FileNotFoundError:
      print(f"metadata file {file_path} does not exist")
      logger.debug(f"metadata file {file_path} does not exist")
      return {}
    
    except yaml.YAMLError:
      print( f"metadata file {file_path} present but does not contain valid yaml")
      logger.error(
        f"metadata file {file_path} present but does not contain valid yaml"
      )
      return {}
    
    return {}

  def predict(self, X, features_names, meta):
    print("GOT A REQ")
    print(meta)
    print(X)
    r = random.randint(0, 10)
    res = np.array([0])
    if (r > 9):
      res = np.array([1])

    return SeldonResponse(data=res)

  def tags(self):
    tag = {
      "model_uri": self.model_uri,
      "model_version": self.mdata["versions"][0]
      }
    return tag
    
  def _create_minio_client(self):
        # Adding prefixing "http" in urlparse is necessary for it to be the netloc
        url = urlparse(os.getenv("AWS_ENDPOINT_URL", "http://s3.amazonaws.com"))
        use_ssl = (
            url.scheme == "https"
            if url.scheme
            else bool(getenv("USE_SSL", "S3_USE_HTTPS", "false"))
        )
        return Minio(
            url.netloc,
            access_key=os.getenv("AWS_ACCESS_KEY_ID", ""),
            secret_key=os.getenv("AWS_SECRET_ACCESS_KEY", ""),
            region=os.getenv("AWS_REGION", ""),
            secure=use_ssl,
        )
  def _download_model(self, uri, temp_dir: str = None):
    client = self._create_minio_client()
    bucket_args = uri.replace("s3://", "", 1).split("/", 1)
    bucket_name = bucket_args[0]
    bucket_path = bucket_args[1] if len(bucket_args) > 1 else ""
    objects = client.list_objects(bucket_name, prefix=bucket_path, recursive=True)
    count = 0
    for obj in objects:
        # Replace any prefix from the object key with temp_dir
        subdir_object_key = obj.object_name.replace(bucket_path, "", 1).strip("/")
        # fget_object handles directory creation if does not exist
        if not obj.is_dir:
            if subdir_object_key == "":
                subdir_object_key = obj.object_name
            client.fget_object(
                bucket_name,
                obj.object_name,
                os.path.join(temp_dir, subdir_object_key),
            )
        count = count + 1
    if count == 0:
        raise RuntimeError(
            "Failed to fetch model. \
            The path or model %s does not exist."
            % (uri)
        )
    return temp_dir
