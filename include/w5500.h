#include "mbed.h"
#include "WIZnetInterface.h"

#define DHCP 1

#define NUCLEO_F411RE

const char *IP_Addr = "192.168.1.2";
const char *IP_Subnet = "255.255.255.0";
const char *IP_Gateway = "192.168.1.1";
unsigned char MAC_Addr[6] = {0x00, 0x08, 0xDC, 0x12, 0x34, 0x56};

void initEthernet(void);
void sendRequest(void);
