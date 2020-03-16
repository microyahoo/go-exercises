#! /usr/bin/env python

def makeActions():
    acts = []
    for i in range(5): # Tries to remember each i
        acts.append(lambda x: i ** x) # All remember same last i!
    return acts

if __name__ == "__main__":
    print makeActions()
    i = 5
    myList = [i, i, i]
    i = 6
    print(myList)
    myList[0] = 7
    print(myList)
