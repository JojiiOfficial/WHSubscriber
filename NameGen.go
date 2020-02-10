package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func getRandomName() string {
	str := "Lorem ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et dolore magna aliqua Non tellus orci ac auctor augue mauris Et netus et malesuada fames ac turpis In hendrerit gravida rutrum quisque In vitae turpis massa sed elementum Odio ut sem nulla pharetra diam sit amet Sit amet aliquam id diam maecenas ultricies Lectus magna fringilla urna porttitor rhoncus dolor Purus in massa tempor nec feugiat nisl pretium fusce id Magnis dis parturient montes nascetur ridiculus mus mauris vitae ultricies Turpis nunc eget lorem dolor sed viverra ipsum nunc Enim ut tellus elementum sagittis vitae et Imperdiet sed euismod nisi porta lorem mollis Sed id semper risus in hendrerit gravida rutrum quisque Vivamus arcu felis bibendum ut tristique et Vivamus at augue eget arcu dictum varius Nec feugiat in fermentum posuere urna Et magnis dis parturient montes A diam maecenas sed enim ut sem viverra aliquet eget Massa massa ultricies mi quis Tellus in metus vulputate eu Fermentum dui faucibus in ornare In est ante in nibh mauris cursus mattis molestie Lacus laoreet non curabitur gravida arcu ac tortor dignissim convallis Quis blandit turpis cursus in Sem nulla pharetra diam sit amet nisl suscipit adipiscing Lectus magna fringilla urna porttitor rhoncus Porttitor leo a diam sollicitudin Leo urna molestie at elementum eu facilisis sed odio Fermentum posuere urna nec tincidunt praesent semper feugiat nibh sed Aliquam eleifend mi in nulla posuere sollicitudin aliquam ultrices sagittis Est lorem ipsum dolor sit amet consectetur Id aliquet risus feugiat in ante metus dictum at tempor A condimentum vitae sapien pellentesque habitant morbi tristique senectus et Arcu dictum varius duis at consectetur lorem donec Sem viverra aliquet eget sit amet Imperdiet nulla malesuada pellentesque elit eget Faucibus a pellentesque sit amet porttitor Sed libero enim sed faucibus turpis in eu Massa tempor nec feugiat nisl pretium Et pharetra pharetra massa massa ultricies Vivamus arcu felis bibendum ut tristique et Nec tincidunt praesent semper feugiat nibh sed Senectus et netus et malesuada fames ac turpis egestas Diam sit amet nisl suscipit adipiscing bibendum Magna sit amet purus gravida Vivamus arcu felis bibendum ut tristique et egestas quis ipsum"
	ss := strings.Split(str, " ")
	r := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(r).Int31n((int32)(len(ss) - 1))
	r2 := rand.New(r).Int31n((int32)(len(ss) - 1))
	return fmt.Sprintf("%s_%s", ss[r1], ss[r2])
}
