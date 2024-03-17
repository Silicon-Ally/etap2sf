package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Silicon-Ally/etap2sf/etap/attachments/exportfiles"
)

const instructions = `

You're likely getting this message because eTap authentication failed. 

In order to retrieve all of your attachments/files, you'll need to authenticate to them as you would through a browser.
You'll need to do this every ~2 hours that you spend downloading attachments, so you'll unfortunately have to do this frequently if you have many files.
For this reason, if you find these instructions don't work well for you, make sure to document your own process for obtaining the authentication information.

Open the steps/step03_download_attachments_from_etap/main.go file in your text editor. There are five values needed for authentication. They are:
	- JSessionID
	- UserDataSessionID
	- AuthSVCToken
	- SecurityToken
	- MyEntityRoleRef
Instructions for each is below. Once you have found one, paste it into the main.go file in the package, and run it.

To find the MyEntityRoleRef:
1 - Navigate to the User entry for you in eTapestry.
2 - In the network tab, the ID will be in the URL of the request to the user's page.
3 - It will look something like 489.0.12345678

To find JSessionID, UserDataSessionID, and AuthSVCToken:
1 - Log into eTapestry in your browser as you normally would.
2 - Navigate to a journal entry that has an attachment on it.
3 - Open the Chrome Dev Tools (F12 or right click -> Inspect).
4 - Go to the Cookies tab, and copy the three cookie's value into the authentication block in the main.go file.

To find the SecurityToken:
1 - On the same page as step 2 from the last set of instructions, open the javascript console in the developer tools.
2 - Run this snippet to find the security token:  

[...document.getElementsByTagName("input")].filter(i => i.name === "securityToken")[0].value

Once you have all five of these values populated in the Authn Block, try to rerun this. If you're still having trouble, please file a bug. 

`

func main() {
	authn := &exportfiles.Authn{
		JSessionID:        "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		UserDataSessionID: "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXXXXXXXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
		AuthSVCToken:      "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		SecurityToken:     "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
		MyEntityRoleRef:   "XXX.X.XXXXXXX",
	}

	if err := exportfiles.Run(authn); err != nil {
		fmt.Print(instructions)
		log.Fatal(err)
	}
	os.Exit(0)
}
