package main

import "log"

const apexClassForEditButton = `

NOTE: IF YOU ARE DEPLOYING INTO PRODUCTION, YOU WILL NEED TO IMPORT THESE FROM A CHANGE SET INSTEAD.
TO DO THAT, CREATE A CHANGE SET FROM THE SANDBOX THAT INCLUDES THE THREE COMPONENTS BELOW, 
PUBLISH IT VIA OUTBOUND CHANGE SET, AND THEN IMPORT IT INTO PRODUCTION VIA INBOUND CHANGE SET.

Salesforce doesn't allow API modification of APEX classes, so you'll need to create a class called EditButton manually.

1. Create a new VisualForcePage for testing:

Setup > Custom Code > Visualforce Pages > New

Name: EditButtonTestComponent
Code:

<apex:page standardController="Account" extensions="EditButton">  
    <p>
        This tab includes data migrated automatically from eTapestry. Editing this data directly is not advised,
        since these fields dont have integrations with other fields in Salesforce.
    </p>
    <apex:form id="etap-migration-form">
        <apex:commandButton value="{!IF(isEditing, 'Cancel', 'Edit')}" action="{!toggleIsEditing}" rerender="etap-migration-form" />
        <apex:commandButton value="Save" action="{!save}" rendered="{!isEditing}" rerender="etap-migration-form" />
        <apex:pageBlock rendered="{!isEditing}">
            <h1>You are not in editing mode</h1>
        </apex:pageBlock>
        <apex:pageBlock rendered="{!NOT(isEditing)}">
            <h1>Note - you are not in editing mode, empty fields are omitted.</h1>
        </apex:pageBlock>
    </apex:form>
</apex:page>


2. Create the Edit Button Apex Class:

Go to: Setup > Apex Classes > New, and paste the following code into the editor:

public class EditButton {
    public Boolean isEditing {get; set;}
    public EditButton (ApexPages.StandardController controller) {
        String modeStr = ApexPages.currentPage().getParameters().get('Mode');
        if(modeStr == 'Edit') {
            isEditing = true;
        } else {
            isEditing = false;
        }
    }
    public void toggleIsEditing() {
        isEditing = !isEditing;
    }
}

3. Create another class called EditButtonTest (you'll need this when importing to production):

@isTest
private class EditButtonTest {
    @isTest static void testConstructorEditMode() {
        // Mock a page and set parameters
        PageReference pageRef = Page.EditButtonTestComponent;
        Test.setCurrentPage(pageRef);
        ApexPages.currentPage().getParameters().put('Mode', 'Edit');

        // Create a new controller instance
        ApexPages.StandardController stdController = new ApexPages.StandardController(new Account()); // Replace 'YourObject' with the actual object
        EditButton controller = new EditButton(stdController);

        // Assert that isEditing is true
        System.assertEquals(true, controller.isEditing, 'isEditing should be true when Mode is Edit');
    }

    @isTest static void testConstructorNonEditMode() {
        PageReference pageRef = Page.EditButtonTestComponent;
        Test.setCurrentPage(pageRef);

        // Create a new controller instance
        ApexPages.StandardController stdController = new ApexPages.StandardController(new Account()); // Replace 'YourObject' with the actual object
        EditButton controller = new EditButton(stdController);

        // Assert that isEditing is false
        System.assertEquals(false, controller.isEditing, 'isEditing should be false when Mode is not Edit');
    }

    @isTest static void testToggleIsEditing() {
        // Mock a page without parameters
        PageReference pageRef = Page.EditButtonTestComponent;
        Test.setCurrentPage(pageRef);

        // Create a new controller instance
        ApexPages.StandardController stdController = new ApexPages.StandardController(new Account()); // Replace 'YourObject' with the actual object
        EditButton controller = new EditButton(stdController);

        // Toggle isEditing and assert changes
        Boolean initialEditMode = controller.isEditing;
        controller.toggleIsEditing();
        System.assertNotEquals(initialEditMode, controller.isEditing, 'isEditing should toggle its value');
    }
}

Once you've created these Apex Classes, you can proceed to the next step.

`

func main() {
	log.Fatalf(apexClassForEditButton)
}
