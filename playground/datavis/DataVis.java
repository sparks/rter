import java.io.*;
import java.util.*;
import processing.core.*;

public class DataVis extends PApplet {
	boolean debug = true;
	
	PFont font;
	int fontSize = 16;
	
	ArrayList<CarAccident> carAccidents;
	
	ScrollingPlot temperaturePlot, humidityPlot;

	public void setup() {
		size(1024, 768);
		background(0);
		
		font = createFont("Arial", fontSize, true);
		
		temperaturePlot = new ScrollingPlot(this, new PVector(100, 100), "Temp.", -100, 100, color(0, 0, 255), color(255, 0, 0));
		humidityPlot = new ScrollingPlot(this, new PVector(300, 100), "Humid.", 0, 100, color(255, 128, 0), color(0, 0, 255));
		
		carAccidents = CarAccident.ParseCsv(this);
	}

	public void draw() {
		
		background(0);
	
		temperaturePlot.update();
		temperaturePlot.draw();
		
		humidityPlot.update();
		humidityPlot.draw();
		
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
