package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
	os.Exit(0)
}

func run() error {
	return fmt.Errorf(`

MANUAL WORK NEEDED:
The following can't be performed via the API, so you'll need to do each manually. Sorry!

MANUAL INTERVETION 1: Adding Visualforce Pages to Layouts

You'll need to manually add the following Visualforce 
components to the following pages:

- Account Soft Credit
- ETap Additional Context
- GAU Allocation
- Partial Soft Credit
- Task

On each of these pages do the following:
	0. Navigate to the Primary Page Layout
	   Setup > Object Manager > [Name of Object Here] > Page Layouts > [Layout Name]
	1. In the pallete at the top of the page, in the left hand side scroll down until you see "Visualforce Pages"
	   Click this, there is a new visualforce page called "eTapestry Migrated Data [Name of Object Here]".
	2. By dragging the "Section" button from the pallete, create a new section at the bottom of the page (typically below "Custom Links" or "System Information"). Select a 1 column layout and name the section "ETapesty" or something similar.
	3. Drag the Visualforce Page from the pallete down into your newly created section.
	4. Click the little wrench icon on the visualforce page you just added, set the height to 1000px, and check the "show scrollbars" box.
	5. Press save (from the pallete), validating that the eTapestry Migrated Data component is at the bottom of the page.

Do this for each page listed above, then go to the next task.


MANUAL INTERVENTION 2: Add File Buttons
NOTE: this is a different (but very similar) task over a DIFFERENT set of objects.

- Account Soft Credit
- ETap Additional Context
- Partial Soft Credit
- Payment
- Task

For each of these objects, do the following:
	0. Navigate to the Primary Page Layout 
		Setup > Object Manager > [Name of Object Here] > Page Layouts > [Layout Name]
	1. In the pallete at the top of the page, in the left hand side scroll down until you see "Related Lists"
	2. Drag the "Files" related list into the page layout.
	3. Press Save. You will be promted to make this change for everyone, click "YES".

MANUAL INTERVENTION 3: Create Content Version Page Layout

0. Navigate to the Page Layouts page for Content Version 
	Setup > Object Manager > Content Version > Page Layouts
1. Click "New", give it a name like "Content Version Layout"
2. Drag "ETap Additional Context" into the page layout somehwere.
3. Save the page.

MANUAL INTERVENTION 4: Add a Related List for Opportunities

0. Navigate to the Page Layouts page for Opportunity
	Setup > Object Manager > Opportunity > Page Layouts
1. For EACH of the layouts on this page:
2. Add an "Account Soft Credit" to the "related list" section at the bottom via drag and drop (in the pallete this is Related Lists > Account Soft Credits).
3. Repeat for each layout.

MANUAL INTERVENTION 5: Update the sort order for activity settings.

0. Navigate to the Activity Settings page
	Setup > Feature Settings > Sales > Activity Settings
1. Uncheck "Sort past activities by the completed date", and press "Submit".

Once you've done these manual interventions, you're done! You can now run the next step.
	`)
}
