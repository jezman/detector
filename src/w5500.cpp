#include "w5500.h"

SPI spi(PA_7, PA_6, PA_5);

#ifdef NUCLEO_F411RE
WIZnetInterface ethernet(&spi, PB_6, PA_9);
#endif

#ifdef BLUE_PILL
WIZnetInterface ethernet(&spi, PB_6, PA_9);
#endif

const char *IP_Addr = "192.168.1.2";
const char *IP_Subnet = "255.255.255.0";
const char *IP_Gateway = "192.168.1.1";
unsigned char MAC_Addr[6] = {0x00, 0x08, 0xDC, 0x12, 0x34, 0x56};

void ethernetUp(void)
{
#if DHCP
    int ret = ethernet.init(MAC_Addr);
#else
    int ret = ethernet.init(MAC_Addr, IP_Addr, IP_Subnet, IP_Gateway);
#endif

    if (!ret)
    {
        printf("Initialized, MAC: %s\r\n", ethernet.getMACAddress());

        ret = ethernet.connect();

        if (!ret)
        {
            printf("IP: %s, MASK: %s, GW: %s\r\n",
                   ethernet.getIPAddress(),
                   ethernet.getNetworkMask(),
                   ethernet.getGateway());
        }
        else
        {
            printf("Error ethernet.connect() - ret = %d\r\n", ret);
            exit(0);
        }
    }
    else
    {
        printf("Error ethernet.init() - ret = %d\r\n", ret);
        exit(0);
    }
}

void ethernetDown(void){
    ethernet.disconnect();
}

void sendRequest(void)
{
    TCPSocketConnection socket;
    socket.connect("www.httpbin.org", 80);

    char sbuffer[] = "GET /anything?a=b HTTP/1.1\r\nHost: www.httpbin.org\r\n\r\n";
    int scount = socket.send(sbuffer, sizeof sbuffer);
    printf("sent %d [%.*s]\r\n", scount, strstr(sbuffer, "\r\n") - sbuffer, sbuffer);

    char rbuffer[64];
    int rcount = socket.receive(rbuffer, sizeof rbuffer);
    printf("recv %d [%.*s]\r\n", rcount, strstr(rbuffer, "\r\n") - rbuffer, rbuffer);

    socket.close();
}

void sendAlert(void)
{
    ethernetUp();
    sendRequest();
    ethernetDown();
}
