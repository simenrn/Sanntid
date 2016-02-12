#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <unistd.h>

int main(){
	int clientSocket, nBytes;
	char buffer[1024];
	struct sockaddr_in clientAddr;
	struct sockaddr_in serverAddr;

	clientSocket = socket(AF_INET, SOCK_DGRAM, 0);
	sendingSocket = socket(AF_INET, SOCK_DGRAM,0);

	memset((char *)&clientAddr, 0, sizeof(clientAddr));
	clientAddr.sin_family = AF_INET;
	clientAddr.sin_port = htons(30000);
	clientAddr.sin_addr.s_addr = htonl(INADDR_ANY);

	memset((char *)&serverAddr, 0, sizeof(serverAddr));
	serverAddr.sin_family = AF_INET;
	serverAddr.sin_port = htons(20021);
	serverAddr.sin_addr.s_addr = inet_addr("129.241.187.23");
	
	if (bind(clientSocket, (struct sockaddr *)&clientAddr, sizeof(clientAddr)) <0){
		printf("bind failed");
		return 0;
	}
	if (bind(serverSocket, (struct sockaddr *)&serverAddr, sizeof(serverAddr)) <0){
		printf("bind failed");
		return 0;
	}
	nBytes = strlen(buffer)+1;
	
	char *my_message = "this is a test message";
	//serverAddr.sin_addr.s_addr = htonl(129.241.187.255)
	while(1){
		if(sendto(clientSocket, my_message, strlen(my_message),0,
		nBytes = recvfrom(clientSocket,buffer,1024,0,NULL,NULL);
		printf("mottatt fra server: %s\n",buffer);
		sleep(4);
	}
	close(clientSocket);
	return 0;
}


//server-ip: 129.241.187.23/44
