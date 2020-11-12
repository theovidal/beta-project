import math

S = [
  [[1, 1], 0],
  [[1, 3], 0],
  [[2, 1], 0],
  [[2, 2], 0],
  [[3, 2], 0],
  [[4, 0.5], 0],
  [[5, 1], 0],

  [[5, 5], 1],
  [[6, 4], 1],
  [[6, 5], 1],
  [[7, 4], 1],
  [[8, 3], 1],
  [[8, 5], 1]
]

Sother = [
  [[1, 1], 0],
  [[2, 2], 0],

  [[4, 4], 1],
  [[1, 4], 1],
  [[5, 5], 1],
]

Snand = [
  [[0, 0], 1],
  [[0, 1], 1],
  [[1, 0], 1],
  [[1, 1], 0],
]

def sigmoid(x):
  return 1.0 / (1.0 + math.exp(x))

def sigmoidPrime(x):
  return sigmoid(x) * (1 - sigmoid(x))

def calculateOutput(x, w, b):
  Σ = b
  for i in range(len(x)):
    Σ += x[i] * w[i]

  if Σ >= 0:
    return 1
  else:
    return 0
  # return sigmoid(Σ)

def deltaWeight(ε, c, o, x):
  return ε * (c - o) * x #* sigmoidPrime(o)

def errorRate(S, w, b):
  Σ = b
  for example in S:
    Σ += (example[1] - calculateOutput(example[0], w, b)) ** 2

  return Σ / 2

def train(ε, epochs, S):
  w = [0, 0]
  b = 0

  for i in range(epochs + 1):
    for example in S:
      x = example[0]
      o = calculateOutput(x, w, b)
      c = example[1]

      b += deltaWeight(ε, c, o, 1)
      for i in range(len(w)):
        w[i] = w[i] + deltaWeight(ε, c, o, x[i])
    print('Epoch {index} --- Error rate: {error}'.format(index=i+1, error=errorRate(S, w, b)))

  return (w, b)


(w, b) = train(0.1, 500, S)
print(w)

ya = -1 * (b / w[1]) / (b / w[0])
yb = (-1 * b) / w[1]

print(f'Line equation: y = {ya}x + {yb}')