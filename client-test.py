from seldon_core.seldon_client import SeldonClient
import numpy as np
import json
import tensorflow_datasets as tfds
import random
import warnings

ds = tfds.load(name="mnist", split="test", as_supervised=True)


sc = SeldonClient(deployment_name="mnist-model", namespace="seldon-system", gateway_endpoint="localhost:8081", gateway="istio")
print(sc.config)
test_size = 10000
corrects = 0
wrong = 0
data = ds.take(test_size).cache()
for image, label in data:

  r = sc.predict(data=np.array(image), gateway="istio",transport="rest")
 # print(r.msg)
  assert(r.success==True)

  res = r.response['data']['tensor']['values']
  print(r.response)
  #print(res)
  prediction = int(np.argmax(np.array(res).squeeze(), axis=0))
  print("predicted: ", prediction, "Truth: ", int(label))
  if int(prediction) == int(label):
    corrects = corrects + 1
  else:
    wrong = wrong + 1
    #print("WRONG", int(prediction), int(label))

print(corrects)
print(wrong)
print(test_size)
print(corrects/test_size)
assert(corrects/test_size > 0.9)