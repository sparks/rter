import java.io.*;
import java.util.*;
import processing.core.*;

public class DataVis extends PApplet {
	boolean debug = true;
	
	PFont font;
	
	ArrayList<CarAccident> carAccidents;
	
	ScrollingPlot temperaturePlot;

	public void setup() {
		size(1024, 768);
		background(0);
		
		font = createFont("Arial", 16, true);
		
		temperaturePlot = new ScrollingPlot(this, new PVector(100, 100), -100, 100, color(0, 0, 255), color(255, 0, 0));
		
		carAccidents = CarAccident.ParseCsv(this);
	}

	public void draw() {
		
		background(0);
	
		temperaturePlot.update();
		temperaturePlot.draw();
		
		/*
		CarAccident test = carAccidents.get(1);
		
		textFont(font, 16);
		fill(255);
		textAlign(LEFT);
		text(test.Latitude.toString(), 10, 20);
		text(test.Longitude.toString(), 10, 60);
		*/
		
		for (CarAccident accident : carAccidents) {
			accident.draw(45.5, -73.6, 1024, 768, 180, 90, 0);
		}
	}

	public void mousePressed() {
		
	}

	public static void main (String [] args) {
		PApplet.main(new String[] { "DataVis" });
	}
}
