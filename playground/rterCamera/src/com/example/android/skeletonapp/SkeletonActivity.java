/*
 * Copyright (C) 2007 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.example.android.skeletonapp;

import android.app.Activity;
import android.content.Intent;
import android.os.Bundle;
import android.provider.Settings.Secure;
import android.view.KeyEvent;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemSelectedListener;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Spinner;
import android.provider.Settings.Secure;

/**
 * This class provides a basic demonstration of how to write an Android
 * activity. Inside of its window, it places a single view: an EditText that
 * displays and edits some internal text.
 */
public class SkeletonActivity extends Activity {
	 
    static final private int BACK_ID = Menu.FIRST;
    static final private int CLEAR_ID = Menu.FIRST + 1;

    private EditText mEditor;
    
    
    int phoneIDselected = 0;
    
    String [] phoneID = {
            "1e7f033bfc7b3625fa07c9a3b6b54d2c81eeff98",
            "fe7f033bfc7b3625fa06c9a3b6b54b2c81eeff98",
            "b6200c5cc15cfbddde2874c40952a7aa25a869dd",
            "852decd1fbc083cf6853e46feebb08622d653602",
            "e1830fcefc3f47647ffa08350348d7e34b142b0b",
            "48ad32292ff86b4148e0f754c2b9b55efad32d1e",
            "acb519f53a55d9dea06efbcc804eda79d305282e",
            "ze7f033bfc7b3625fa06c5a316b54b2c81eeff98",
            "t6200c5cc15cfbddde2875c41952a7aa25a869dd",
            "952decd1fbc083cf6853e56f1ebb08622d653602",
            "y1830fcefc3f47647ffa05351348d7e34b142b0b",
            "x8ad32292ff86b4148e0f55412b9b55efad32d1e",
            "qcb519f53a55d9dea06ef5cc104eda79d305282e" };
    
    public SkeletonActivity() {
    }

    /** Called with the activity is first created. */
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        // Inflate our UI from its XML layout description.
        setContentView(R.layout.skeleton_activity);
        
        // populate spinner
        String [] phoneIDnumber = {
        "phone ID 1",
        "phone ID 2",
        "phone ID 3",
        "phone ID 4",
        "phone ID 5",
        "phone ID 6",
        "phone ID 7",
        "phone ID 8",
        "phone ID 9",
        "phone ID 10",
        "phone ID 11",
        "phone ID 12",
        "phone ID 13" };
        
        Spinner idSelect = (Spinner) findViewById(R.id.phone_id);
        ArrayAdapter<String> idAdapter = new ArrayAdapter<String>(this,
        		android.R.layout.simple_spinner_dropdown_item, phoneIDnumber);
        idSelect.setAdapter(idAdapter);
        
        // Find the text editor view inside the layout, because we
        // want to do various programmatic things with it.
//        mEditor = (EditText) findViewById(R.id.editor);

        // Hook up button presses to the appropriate event handler.
//        ((Button) findViewById(R.id.back)).setOnClickListener(mBackListener);
        ((Button) findViewById(R.id.camera)).setOnClickListener(mCameraListener);
        
        idSelect.setOnItemSelectedListener(mPhoneIDSelectListener);
        
//        mEditor.setText(getText(R.string.main_label));
    }

    /**
     * Called when the activity is about to start interacting with the user.
     */
    @Override
    protected void onResume() {
        super.onResume();
    }

    /**
     * Called when your activity's options menu needs to be created.
     */
    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        super.onCreateOptionsMenu(menu);

        // We are going to create two menus. Note that we assign them
        // unique integer IDs, labels from our string resources, and
        // given them shortcuts.
        menu.add(0, BACK_ID, 0, R.string.back).setShortcut('0', 'b');
        menu.add(0, CLEAR_ID, 0, R.string.clear).setShortcut('1', 'c');

        return true;
    }

    /**
     * Called right before your activity's option menu is displayed.
     */
    @Override
    public boolean onPrepareOptionsMenu(Menu menu) {
        super.onPrepareOptionsMenu(menu);

        // Before showing the menu, we need to decide whether the clear
        // item is enabled depending on whether there is text to clear.
//        menu.findItem(CLEAR_ID).setVisible(mEditor.getText().length() > 0);

        return true;
    }

    /**
     * Called when a menu item is selected.
     */
    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        
    	
    	switch (item.getItemId()) {
        case BACK_ID:
            finish();
            return true;
        case CLEAR_ID:
        	
            return true;
        }

        return super.onOptionsItemSelected(item);
    }

    /**
     * A call-back for when the user presses the back button.
     */
    OnClickListener mBackListener = new OnClickListener() {
        public void onClick(View v) {
            finish();
        }
    };

    /**
     * A call-back for when the user presses the clear button.
     */
    OnClickListener mCameraListener = new OnClickListener() {
        public void onClick(View v) {
        	openCamera();
        }	
    };
    
    /**
     * A call-back for when the user presses the clear button.
     */
    OnItemSelectedListener mPhoneIDSelectListener = new OnItemSelectedListener() {

		@Override
		public void onItemSelected(AdapterView<?> arg0, View arg1, int position,
				long arg3) {
			changePhoneID(position);
		}

		@Override
		public void onNothingSelected(AdapterView<?> arg0) {
			// TODO Auto-generated method stub
			
		}
    };
    
    private void openCamera() {
		// TODO Auto-generated method stub
    	Intent i = new Intent(this, CameraPreview.class);
    	i.putExtra("phoneID", this.phoneID[this.phoneIDselected]);
    	startActivity(i);
	}
    
    private void changePhoneID(int id) {
    	this.phoneIDselected = id;
    }
}
