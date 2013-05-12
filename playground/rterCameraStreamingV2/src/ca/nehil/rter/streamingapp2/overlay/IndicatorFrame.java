package ca.nehil.rter.streamingapp2.overlay;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.FloatBuffer;

import javax.microedition.khronos.opengles.GL10;

public class IndicatorFrame {
	private FloatBuffer vertexBuffer; // Buffer for vertex-array
	
	public static enum Colour {
		RED, GREEN, BLUE
	}

	private float[] vertices = { // Vertices for the arrow
			// TOP
			-1.0f, 0.9f, 0.0f, // 0. left-bottom-top
			1.0f, 0.9f, 0.0f, // 1. right-bottom-top
			-1.0f, 1.0f, 0.0f, // 2. left-top-top
			1.0f, 1.0f, 0.0f, // 3. right-top-top
			// LEFT
			-1.0f, -0.9f, 0.0f, // 0. left-bottom-left
			-0.9f, -0.9f, 0.0f, // 1. right-bottom-left
			-1.0f, 0.9f, 0.0f, // 2. left-top-left
			-0.9f, 0.9f, 0.0f, // 3. right-top-left
			// BOTTOM
			-1.0f, -1.0f, 0.0f, // 0. left-bottom-bottom
			1.0f, -1.0f, 0.0f, // 1. right-bottom-bottom
			-1.0f, -0.9f, 0.0f, // 2. left-top-bottom
			1.0f, -0.9f, 0.0f, // 3. right-top-bottom
			// RIGHT
			0.9f, -0.9f, 0.0f, // 0. left-bottom-left
			1.0f, -0.9f, 0.0f, // 1. right-bottom-left
			0.9f, 0.9f, 0.0f, // 2. left-top-left
			1.0f, 0.9f, 0.0f // 3. right-top-left
	};
	
	private Colour currentColour = Colour.BLUE;

	// Constructor - Setup the vertex buffer
	public IndicatorFrame() {		
		// Setup vertex array buffer. Vertices in float. A float has 4 bytes
		ByteBuffer vbb = ByteBuffer.allocateDirect(vertices.length * 4);
		vbb.order(ByteOrder.nativeOrder()); // Use native byte order
		vertexBuffer = vbb.asFloatBuffer(); // Convert from byte to float
		this.resize(1.0f, 1.0f, 0.0f);
	}
	
	public void resize(float xTotal, float yTotal, float distance) {
		float framePercentWidth = 0.025f;
		float frameWidth;
		//use the smallest dimension to determine frame width
		if (xTotal < yTotal) {
			frameWidth = framePercentWidth*xTotal;
		} else {
			frameWidth = framePercentWidth*yTotal;
		}
		float right = xTotal/2.0f;
		float left = -right;
		float top = yTotal/2.0f;
		float bottom = -top;
		distance = -distance;
		
		
		float[] vertices_tmp =
			{ // Vertices for the arrow
				// TOP
				left, top-frameWidth, distance, // 0. left-bottom-top
				right, top-frameWidth, distance, // 1. right-bottom-top
				left, top, distance, // 2. left-top-top
				right, top, distance, // 3. right-top-top
				// LEFT
				left, bottom+frameWidth, distance, // 0. left-bottom-left
				left+frameWidth, bottom+frameWidth, distance, // 1. right-bottom-left
				left, top-frameWidth, distance, // 2. left-top-left
				left+frameWidth, top-frameWidth, distance, // 3. right-top-left
				// BOTTOM
				left, bottom, distance, // 0. left-bottom-bottom
				right, bottom, distance, // 1. right-bottom-bottom
				left, bottom+frameWidth, distance, // 2. left-top-bottom
				right, bottom+frameWidth, distance, // 3. right-top-bottom
				// RIGHT
				right-frameWidth, bottom+frameWidth, distance, // 0. left-bottom-left
				right, bottom+frameWidth, distance, // 1. right-bottom-left
				right-frameWidth, top-frameWidth, distance, // 2. left-top-left
				right, top-frameWidth, distance // 3. right-top-left
		};
		
		vertices = vertices_tmp;
		vertexBuffer.put(vertices); // Copy data into buffer
		vertexBuffer.position(0); // Rewind
	}

	// Render the shape
	public void draw(GL10 gl) {
		// Enable vertex-array and define its buffer
		gl.glEnableClientState(GL10.GL_VERTEX_ARRAY);
		gl.glVertexPointer(3, GL10.GL_FLOAT, 0, vertexBuffer);

		// Render all the faces
		for (int face = 0; face < 4; face++) {
			// Set the color for each of the faces
			switch (this.currentColour) {
			case RED:
				gl.glColor4f(0.9f, 0.0f, 0.0f, 1.0f); //0.8f);
				break;
			case GREEN:
				gl.glColor4f(0.0f, 0.9f, 0.0f, 1.0f); //0.8f);
				break;
			case BLUE:
				gl.glColor4f(0.0f, 0.0f, 0.9f, 1.0f); //0.8f);
				break;
			}

			// Draw the primitive from the vertex-array directly
			gl.glDrawArrays(GL10.GL_TRIANGLE_STRIP, face * 4, 4);
		}
		gl.glDisableClientState(GL10.GL_VERTEX_ARRAY);
	}
	
	/**
	 * not thread safe
	 * 
	 * @param colour
	 */
	public void colour(Colour colour){
		this.currentColour = colour;
	}
}
