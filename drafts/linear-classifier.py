import math

S = [
  [[1, 1], 0],
  [[1, 3], 0],
  [[2, 1], 0],
  [[2, 2], 0],
  [[3, 2], 0],
  [[4, 0.5], 0],
  [[5, 1], 0],

  [[6.5, 1.75], 1],
  [[5, 5], 1],
  [[6, 4], 1],
  [[6, 5], 1],
  [[7, 4], 1],
  [[8, 3], 1],
  [[8, 5], 1],
  [[17, 12], 1]
]

def σ(x):
  return 1.0 / (1.0 + math.exp(-x))

def σPrime(x):
  return σ(x) * (1 - σ(x))

def calculateOutput(x, w, b):
  Σ = b
  for i in range(len(x)):
    Σ += x[i] * w[i]

  return σ(Σ)

def deltaWeight(ε, c, o, x):
  return ε * (c - o) * x * σPrime(o)

def errorRate(S, w, b):
  Σ = b
  for example in S:
    Σ += (example[1] - calculateOutput(example[0], w, b)) ** 2

  return Σ / 2

def train(ε, epochs, S):
  w = [0, 0]
  b = 0

  epoch = 0
  for i in range(epochs):
    epoch += 1
    for example in S:
      x = example[0]
      o = calculateOutput(x, w, b)
      c = example[1]

      b += deltaWeight(ε, c, o, 1)
      for i in range(len(w)):
        w[i] = w[i] + deltaWeight(ε, c, o, x[i])
    
    print('Epoch {} --- Error rate: {}'.format(epoch, errorRate(S, w, b)))

  return (w, b)


(w, b) = train(0.1, 500, S)
print('')
print('Weights: {}'.format(w))

𝑎 = -1 * (b / w[1]) / (b / w[0])
𝑏 = (-1 * b) / w[1]

print('Line equation: 𝑦 = {}𝑥 + {}'.format(𝑎, 𝑏))
