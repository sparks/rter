import panoia.*;
import java.io.*;
import java.util.*;
import processing.core.*;

public class DataVis extends PApplet {
	boolean debug = true;
	
	PFont font;
	
	ArrayList<CarAccident> carAccidents;

	public void setup() {
		size(displayWidth, displayHeight);
		background(0);

		smooth();
		
		font = createFont("Arial", 16, true);
		
		carAccidents = CarAccident.ParseCsv(this);
	}

	public void draw() {
		background(0);

		/*
		CarAccident test = carAccidents.get(1);
		
		textFont(font, 16);
		fill(255);
		textAlign(LEFT);
		text(test.Latitude.toString(), 10, 20);
		text(test.Longitude.toString(), 10, 60);
		*/
		
		for (CarAccident accident : carAccidents) {
			accident.draw(45.5, -73.6, 1600, 900, 180, 90, 0);
		}
	}

	public void mousePressed() {
		
	}

	public static void main (String [] args) {
		PApplet.main(new String[] { "DataVis" });
	}
}
