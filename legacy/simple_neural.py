import math
import numpy as np

# The examples base, to train our network
S = [
  [np.array([1, 1]), np.array([1, 0])],
  [np.array([1, 3]), np.array([1, 0])],
  [np.array([2, 1]), np.array([1, 0])],
  [np.array([2, 2]), np.array([1, 0])],
  [np.array([3, 2]), np.array([1, 0])],
  [np.array([4, 0.5]), np.array([1, 0])],
  [np.array([5, 1]), np.array([1, 0])],

  [np.array([6.5, 1.75]), np.array([0, 1])],
  [np.array([5, 5]), np.array([0, 1])],
  [np.array([6, 4]), np.array([0, 1])],
  [np.array([6, 5]), np.array([0, 1])],
  [np.array([7, 4]), np.array([0, 1])],
  [np.array([8, 3]), np.array([0, 1])],
  [np.array([8, 5]), np.array([0, 1])],
  [np.array([17, 12]), np.array([0, 1])]
]

def σ(x):
  return 1.0 / (1.0 + math.exp(-x))

def σPrime(x):
  return σ(x) * (1 - σ(x))

σ_v = np.vectorize(σ)
σPrime_v = np.vectorize(σPrime)

class Network:
  # Create a new neural network
  def __init__(self, ε):
    self.ε = ε
    self.weights = [
      np.array([[0, 0, 0],
                [0, 0, 0]], dtype=float),
      np.array([[0, 0],
                [0, 0],
                [0, 0]], dtype=float)
    ]
    self.biases = [
      np.array([0, 0, 0], dtype=float), np.array([0, 0], dtype=float)
    ]
    self.results = [
      np.array([0, 0], dtype=float),
      np.array([0, 0, 0], dtype=float),
      np.array([0, 0], dtype=float)
    ]
    self.outputs = [
      np.array([0, 0], dtype=float),
      np.array([0, 0, 0], dtype=float),
      np.array([0, 0], dtype=float)
    ]

  # Calculate the output of the network
  def feed_forward(self, example):
    self.outputs[0] = example.copy()
    self.results[0] = example.copy()
    for layer in range(len(self.outputs) - 1):
      self.results[layer + 1] = np.dot(self.outputs[layer], self.weights[layer]) + self.biases[layer]
      self.outputs[layer + 1] = σ_v(self.results[layer + 1])

  # Adjust the weights of the neural network
  def feed_backward(self, example):
    δ = [
      np.array([0, 0], dtype=float),
      np.array([0, 0, 0], dtype=float),
      np.array([0, 0], dtype=float)
    ]

    # Calculating derivative for the last layer
    for neuron in range(len(self.outputs[-1])):
      δ[-1][neuron] = σPrime(self.results[-1][neuron]) * (self.outputs[-1][neuron] - example[neuron])

    # Calculating derivatives for the next layers
    for layer in range(len(self.outputs) - 2, -1, -1):
      for neuron in range(len(self.outputs[layer])):
        δ[layer][neuron] = σPrime(self.results[layer][neuron]) * np.dot(self.weights[layer][neuron], np.transpose(δ[layer + 1]))

    # Apply changes to weights
    for layer in range(len(self.weights)):
      for j in range(len(self.weights[layer])):
        for i in range(len(self.weights[layer][j])):
          self.weights[layer][j][i] = self.weights[layer][j][i] + self.ε * δ[layer + 1][i] * self.outputs[layer][j]

  # Train the network with the example base S and do a set number of epochs
  def train(self, S, epochs):
    for epoch in range(epochs):
      for example in S:
        self.feed_forward(example[0])
        self.feed_backward(example[1])

      print('Epoch {}'.format(epoch + 1))

neural = Network(0.1)
neural.train(S, 50)
print(neural.weights)

neural.feed_forward(np.array([3, 1]))
print(neural.outputs[-1])
neural.feed_forward(np.array([8, 5]))
print(neural.outputs[-1])
