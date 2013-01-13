package com.example.android.skeletonapp.overlay;

import javax.microedition.khronos.egl.EGLConfig;
import javax.microedition.khronos.opengles.GL10;
import android.content.Context;
import android.opengl.GLSurfaceView;
import android.opengl.GLU;

import android.opengl.GLSurfaceView.Renderer;

public class CameraGLRenderer implements Renderer {
	Arrow arrowLeft;
	Arrow arrowRight;
	IndicatorFrame indicatorFrame;

	Context context;   // Application's context
	
	float aspect;
	float xTotal, yTotal, distance;
	
	float arrowScaleSpeed = 0.05f;
	float arrowScaleSpeedMin = 0.008f;
	float arrowScale = 1.0f;
	float arrowScaleMax = 1.2f;
	float arrowScaleMin = 1.0f;
	boolean arrowScaleIncrease = true;
	   
	   // Constructor with global application context
	   public CameraGLRenderer(Context context) {
	      this.context = context;
	      
	      arrowLeft = new Arrow();
	      arrowRight = new Arrow();
	      indicatorFrame = new IndicatorFrame();
	   }
	   
	   // Call back when the surface is first created or re-created
	   public void onSurfaceCreated(GL10 gl, EGLConfig config) {
	      gl.glClearColor(0.0f, 0.0f, 0.0f, 0.0f);  // Set color's clear-value to black
	      gl.glClearDepthf(1.0f);            // Set depth's clear-value to farthest
	      gl.glEnable(GL10.GL_DEPTH_TEST);   // Enables depth-buffer for hidden surface removal
	      gl.glDepthFunc(GL10.GL_LEQUAL);    // The type of depth testing to do
	      gl.glHint(GL10.GL_PERSPECTIVE_CORRECTION_HINT, GL10.GL_NICEST);  // nice perspective view
	      gl.glShadeModel(GL10.GL_SMOOTH);   // Enable smooth shading of color
	      gl.glDisable(GL10.GL_DITHER);      // Disable dithering for better performance
	  
	      // You OpenGL|ES initialization code here
	      // ......
	   }
	   
	   // Call back after onSurfaceCreated() or whenever the window's size changes
	   public void onSurfaceChanged(GL10 gl, int width, int height) {
	      if (height == 0) height = 1;   // To prevent divide by zero
	      aspect = (float)width / height;
	      
	      // get the total x and y at distance
	      distance = 6.0f;
	      xTotal = (float) (aspect*Math.tan(Math.toRadians(45.0/2))*distance*2);
	      yTotal = (float) (Math.tan(Math.toRadians(45.0/2))*distance*2);
	      
	      indicatorFrame.resize(xTotal, yTotal, distance);
	   
	      // Set the viewport (display area) to cover the entire window
	      gl.glViewport(0, 0, width, height);
	  
	      // Setup perspective projection, with aspect ratio matches viewport
	      gl.glMatrixMode(GL10.GL_PROJECTION); // Select projection matrix
	      gl.glLoadIdentity();                 // Reset projection matrix
	      // Use perspective projection
	      GLU.gluPerspective(gl, 45, aspect, 0.1f, 100.f);
	  
	      gl.glMatrixMode(GL10.GL_MODELVIEW);  // Select model-view matrix
	      gl.glLoadIdentity();                 // Reset
	  
	      // You OpenGL|ES display re-sizing code here
	      // ......
	      
	   }
	   
	   // Call back to draw the current frame.
	   public void onDrawFrame(GL10 gl) {
	      // Clear color and depth buffers using clear-value set earlier
	      gl.glClear(GL10.GL_COLOR_BUFFER_BIT | GL10.GL_DEPTH_BUFFER_BIT);
	     
	      // You OpenGL|ES rendering code here
	      // ......
	      
	      // pulsate arrows
	      if (arrowScaleIncrease) {
	    	  float speed = (arrowScaleMax - arrowScale)*arrowScaleSpeed;
	    	  if (speed < arrowScaleSpeedMin) speed = arrowScaleSpeedMin;
	    	  arrowScale += speed;
	    	  if (arrowScale >= arrowScaleMax) {
	    		  arrowScale = arrowScaleMax;
	    		  arrowScaleIncrease = false;
	    	  }
	      } else {
	    	  float speed = (arrowScale - arrowScaleMin)*arrowScaleSpeed;
	    	  if (speed < arrowScaleSpeedMin) speed = arrowScaleSpeedMin;
	    	  arrowScale -= speed;
	    	  if (arrowScale <= arrowScaleMin) {
	    		  arrowScale = arrowScaleMin;
	    		  arrowScaleIncrease = true;
	    	  }
	      }
	      
	      // FRAME
	      gl.glLoadIdentity();
	      indicatorFrame.draw(gl);
	      
	      // RIGHT ARROW
	      gl.glLoadIdentity();                 // Reset model-view matrix ( NEW )
	      gl.glTranslatef(xTotal/2.0f - 15.0f/xTotal, 0.0f, -distance); // Translate left and into the screen ( NEW )
	      gl.glScalef(arrowScale, arrowScale, 1.0f);
	      arrowRight.draw(gl);                   // Draw triangle ( NEW )
	  
	      // LEFT
	      gl.glLoadIdentity();
	      gl.glTranslatef(-xTotal/2.0f + 15.0f/xTotal, 0.0f, -distance); // Translate left and into the screen ( NEW )
	      gl.glRotatef(180.0f, 0.0f, 0.0f, 1.0f);
	      gl.glScalef(arrowScale, arrowScale, 1.0f);
	      arrowLeft.draw(gl);                       // Draw quad ( NEW )
	   }

}
