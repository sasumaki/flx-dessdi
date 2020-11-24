
from cloudevents.http import from_http
from typing import List, Dict
from ceserver import CEServer 
from ceserver import CEModel
from enum import Enum
import tempfile
import os
from urllib.parse import urlparse
from minio import Minio
import onnxruntime as rt


class Protocol(Enum):
    tensorflow_http = "tensorflow.http"
    seldon_http = "seldon.http"

    def __str__(self):
        return self.value
class DummyModel(CEModel):

    def __init__(self, name:str, model_uri:str = None, create_response:bool = True):
        super().__init__(name)
        self.create_response = create_response
        self.model_uri = model_uri

    @staticmethod
    def getResponse() -> Dict:
        return {"foo": 1}

    def load(self):
      
      print("wat")
      self.out_dir = tempfile.mkdtemp()

      model_file =  os.path.join(self._download_model(self.model_uri, self.out_dir), "model.onnx")

      self._model = model_file
      print(self._model)
      self.session = rt.InferenceSession(self._model, None)
      self.input_name = self.session.get_inputs()[0].name
      self.output_name = self.session.get_outputs()[0].name
      self.ready = True
      print("init and model loading done!!!")
      return

    def process_event(self, inputs: List, headers: Dict) -> Dict:
        assert headers[customHeaderKey] == customHeaderVal
        if self.create_response:
            return DummyModel.getResponse()
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


if __name__ == "__main__":
  server = CEServer(
          Protocol.seldon_http, 9000
      )
  uri =  os.getenv("MODEL_URI")

  model = DummyModel("name", uri)
  model.load()

  server.start(model)
 
 