static const char wifi_ssid[] = "homewifi";
static const char wifi_pass[] = "password";

/* Yout must change XOR key before prod setup */
static const char xor_key[]   = "10101010";

/* Controller server informations */
static char path[] = "/smart-home";
static char host[] = "192.168.1.69";
static int port = 8089;

/* Debug logging */
static const int debug = 0; /* 0 mens off, 1 means on */