#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <unistd.h>

int main(){
	int clientSocket, nBytes;
	char buffer[1024];
	struct sockaddr_in serverAddr;

	clientSocket = socket(AF_INET, SOCK_DGRAM, 0);

	memset((char *)&serverAddr, 0, sizeof(serverAddr));
	serverAddr.sin_family = AF_INET;
	serverAddr.sin_port = htons(30000);
	serverAddr.sin_addr.s_addr = htonl(INADDR_ANY);
	
	if (bind(clientSocket, (struct sockaddr *)&serverAddr, sizeof(serverAddr)) <0){
		printf("bind failed");
		return 0;
	}
	printf("bind completed");
	nBytes = strlen(buffer)+1;
	
	
	//serverAddr.sin_addr.s_addr = htonl(129.241.187.255)
	while(1){
		nBytes = recvfrom(clientSocket,buffer,1024,0,NULL,NULL);
		printf("mottatt fra server: %s\n",buffer);
		sleep(4);
	}
	close(clientSocket);
	return 0;
}
