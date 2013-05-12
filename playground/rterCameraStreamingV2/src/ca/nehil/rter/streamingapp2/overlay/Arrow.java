package ca.nehil.rter.streamingapp2.overlay;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.FloatBuffer;
import javax.microedition.khronos.opengles.GL10;

/**
 * right facing arrow
 * the origin is on the left edge of the arrow
 * @author stepan
 *
 */
public class Arrow {
	private FloatBuffer vertexBuffer; // Buffer for vertex-array

	private float[] vertices = { // Vertices for the arrow
			-0.5f, -1.0f, 0.0f, // 0. left-bottom
			0.5f, 0.0f, 0.0f, // 1. right-bottom
			-0.5f, 1.0f, 0.0f, // 2. left-top
	};

	// Constructor - Setup the vertex buffer
	public Arrow() {
		// Setup vertex array buffer. Vertices in float. A float has 4 bytes
		ByteBuffer vbb = ByteBuffer.allocateDirect(vertices.length * 4);
		vbb.order(ByteOrder.nativeOrder()); // Use native byte order
		vertexBuffer = vbb.asFloatBuffer(); // Convert from byte to float
		vertexBuffer.put(vertices); // Copy data into buffer
		vertexBuffer.position(0); // Rewind
	}

	// Render the shape
	public void draw(GL10 gl) {
		// Enable vertex-array and define its buffer
		gl.glEnableClientState(GL10.GL_VERTEX_ARRAY);
		gl.glVertexPointer(3, GL10.GL_FLOAT, 0, vertexBuffer);
		gl.glColor4f(0.9f, 0.0f, 0.0f, 1.0f); //0.65f);      // Set the current color (NEW)
		// Draw the primitives from the vertex-array directly
		gl.glDrawArrays(GL10.GL_TRIANGLE_STRIP, 0, vertices.length / 3);
		gl.glDisableClientState(GL10.GL_VERTEX_ARRAY);
	}
}
