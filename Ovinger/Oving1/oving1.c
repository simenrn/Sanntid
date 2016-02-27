#include <pthread.h>
#include <stdio.h>

int i = 0;
int up = 0;
int down = 0;
void* increment(){
	for (up = 0; up < 1000000; up+=1){
		i+=1;
	}
	return NULL;
}

void* decrement(){
	for (down = 0; down < 1000000; down+=1){
		i-=1;
	}
	return NULL;
}

int main(){
	pthread_t inc;
	pthread_t dec;

	pthread_create(&inc, NULL, increment, NULL);
	pthread_create(&dec, NULL, decrement, NULL);

	pthread_join(inc, NULL);
	pthread_join(dec, NULL);
	printf("i: %d\n", i);
	return 0;
}
