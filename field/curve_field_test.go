package field

import (
	"testing"
)

/*
    const char *genInitStr = "[7852334875614213225969535005319230321249629225894318783946607976937179571030765324627135523985138174020408497250901949150717492683934959664497943409406486,8189589736511278424487290408486860952887816120897672059241649987466710766123126805204101070682864313793496226965335026128263318306025907120292056643404206]";

    ProxyReEncrypt afgh = ProxyReEncrypt::initFromFile("../pbc-0.5.14/a.param", genInitStr);

    element_t secretKey1;
    afgh.createSecretKey(secretKey1, "276146606970621369032156664792541580771690346936");
    // dumpKeyToConsole(secretKey1);

    element_t secretKey2;
    afgh.createSecretKey(secretKey2, "723726390432394136624550768129737051661740488013");
    // dumpKeyToConsole(secretKey2);

    element_t publicKey1;
    afgh.generatePublicKey(secretKey1, publicKey1);
    checkPointElement(publicKey1,
                      "2280014378744220144146373205831932526719685024545487661471834655738123196933971699437542834115250416780965121862860444719075976277314039181516434962834201",
                      "5095219617050150661236292739445977344231341334112418835906977843435788249723740037212827151314561006651269984991436149205169409973600265455370653168534480");

    element_t publicKey2;
    afgh.generatePublicKey(secretKey2, publicKey2);
    checkPointElement(publicKey2,
                      "3179305015224135534600913697529322474159056835318317733023669075672623777446135077445204913153064372784513253897383442505556687490792970594174506652914922",
                      "6224780198151226873973489249032341465791104108959675858554195681300248102506693704394812888080668305563719500114924024081270783773410808141172779117133345");

 */

func TestCurveElements(t *testing.T) {
}