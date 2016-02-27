
from threading import Thread

i=0



def thread_1():
        global i
        for p in range(1000000):
                i += 1
	

def thread_2():
        global i
        for p in range(1000000):
                i -= 1
	



def main():
	inc = Thread(target = thread_1, args = (),)
	dec = Thread(target = thread_2, args = (),)
	inc.start()
	dec.start()
	inc.join()
	dec.join()
	print(i)
	
main()
