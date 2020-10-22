from seldon_core.seldon_client import SeldonClient
import numpy as np
import matplotlib.pyplot as plt
import json
from tensorflow.examples.tutorials.mnist import input_data
import random

mnist = input_data.read_data_sets("MNIST_data/", one_hot=False)

X_train = np.vstack([img.reshape(-1,) for img in mnist.train.images])
y_train = mnist.train.labels

X_test = np.vstack([img.reshape(-1,) for img in mnist.test.images])
y_test = mnist.test.labels

sc = SeldonClient(deployment_name="mnist-model",namespace="seldon-system",gateway_endpoint="localhost:8081")
print(sc)

data = random.choice(X_test)
plt.imshow(data.reshape(28,28))
plt.colorbar()
plt.show()
r = sc.predict(data=data, gateway="istio",transport="rest")
assert(r.success==True)

res = r.response['data']['tensor']['values']
print(res)

print(int(np.argmax(np.array(res).squeeze(), axis=0)))