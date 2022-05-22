#include <stdio.h>
#include <sys/socket.h>
#include <linux/if_packet.h>
#include <net/ethernet.h>


int main() {
    int packet_socket;
    packet_socket = socket(PF_PACKET, SOCK_RAW, htons(ETH_P_ALL));
    if (packet_socket < 0) {
        perror("socket");
        return 1;
    }
}

