import java.io.*;
import java.util.*;
import processing.core.*;

public class ScrollingPlot {
	
	final int barCount = 10;
	final float graphSize;
	
	PApplet parent;
	PVector position;
	String title;
	
	float min, max, avg;
	int low, high, mid;
	
	int time, oldTime;
	float[] values;

	PFont font;
	int fontSize = 16;
	
	public ScrollingPlot(PApplet parent, PVector position, String title, float min, float max, int low, int high, int graphSize) {
		this.parent = parent;
		this.position = position;
		this.title = title;
		this.min = min;
		this.max = max;
		this.avg = (min + max) / 2;
		this.low = low;
		this.high = high;
		this.mid = parent.color(255);
		this.graphSize = graphSize;
		
		values = new float[barCount];

		font = parent.createFont("Arial", fontSize, true);
		
		for (int i = 0; i < barCount; i++)
			values[i] = avg;
		
		time = oldTime = parent.second();
	}
	
	public void update() {
		time = parent.second();
		
		// The below simulates real-time data being received at a nice rate of 1Hz, for now.
		if (time != oldTime) {
			for (int i = 0; i < barCount - 1; i++) {
				values[i] = values[i + 1];
			}
			
			values[barCount - 1] = parent.random(min, max);
		}
		
		oldTime = time;
	}
	
	public void draw() {
		
		float middle = position.y + graphSize / 2;
		
		// Fill plot with pretty gradients.
		for (int i = 0; i < barCount; i++) {
			if (values[i] > avg) {
			
				if (values[i] > max) values[i] = max;
				
				float amount = (values[i] - avg) / (max - avg);
				int top = (int)(position.y + (1 - amount) * graphSize / 2);
				int height = (int)(amount * graphSize / 2);
				int c = parent.lerpColor(mid, high, amount);
				
				setGradient((int)position.x + (int)(i * graphSize / barCount), top, graphSize / barCount, height, mid, c);
			} else {
			
				if (values[i] < min) values[i] = min;
				
				float amount = (values[i] - avg) / (min - avg);
				int top = (int)middle;
				int height = (int)(amount * graphSize / 2);
				int c = parent.lerpColor(mid, low, amount);
				
				setGradient((int)position.x + (int)(i * graphSize / barCount), top, graphSize / barCount, height, c, mid);
			}
		}
		
		// Draw plot axis lines.
		parent.stroke(255);
		parent.fill(255);
		parent.line(position.x, position.y, position.x, position.y + graphSize);
		parent.line(position.x, middle, position.x + graphSize, middle);
		
		// Print plot title and y-axis values.
		parent.textFont(font, fontSize);
		parent.textAlign(PConstants.CENTER);
		parent.text(title, position.x + graphSize / 2, position.y - 10);
		parent.textAlign(PConstants.RIGHT);
		parent.text(String.format("%.1f", max), position.x - 5, position.y + fontSize);
		parent.text(String.format("%.1f", min), position.x - 5, position.y + graphSize);
	}
	
	private void setGradient(int x, int y, float w, float h, int c1, int c2) {

		parent.noFill();

		for (int i = y; i <= y + h; i++) {
			float inter = parent.map(i, y, y + h, 0, 1);
			int c = parent.lerpColor(c2, c1, inter);
			parent.stroke(c);
			parent.line(x, i, x + w, i);
		}
	}
	
}
