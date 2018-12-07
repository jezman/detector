#include "w5500.h"

DigitalIn detector_1(PA_8);
DigitalIn detector_2(PB_5);
DigitalIn detector_3(PA_10);

int main()
{

    if (detector_1 || detector_2 || detector_3)
    {
        ethernetUp();
        sendRequest();
        ethernetDown();
    }

    return 0;
}
