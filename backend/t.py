def f(a, b):
    if b == 0:
        return a
    else:
        return f(a +1, b -1)

def g(a, b):
    if b == 0:
        return 0
    else:
        return f(a, g(a, b -1))

x = g(6, 4)
print(x)