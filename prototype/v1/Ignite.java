import panoia.*;
import processing.core.*;
import java.io.*;
import java.util.*;
import java.lang.Math;

public class Ignite extends PApplet {

	Pano pano;
	PanoPov pov;

	float fov = 360;

	ArrayList<CarAccident> carAccidents;
	ArrayList<BikeAccident> bikeAccidents;

	PFont font;

	public void setup() {
		size(displayWidth, (int)(displayWidth/(6.5f*fov/360)));
		// size(1024*3, 768);

		background(0);

		font = createFont("Helvetica", 14);
		textFont(font, 14);

		pano = new Pano(this);
		pov = pano.getPov();
		// pano.setPano("3qry8ACTZ8Mw6SQ1UaLNMg");
		pano.setPosition(new LatLng(45.5110809f, -73.5700496f));
		// pano.setPosition(new LatLng(45.52059937f, -73.58165741f));


		LatLng a = new LatLng(45.537719f, -73.622271f); //3qry8ACTZ8Mw6SQ1UaLNMg
		LatLng b = new LatLng(45.537558f, -73.621739f); //gUV-eHyqZN3NlJQA6pZQNw

		smooth();
		stroke(255);
		fill(255);
		// noLoop();

		carAccidents = CarAccident.ParseCsv(this);
		bikeAccidents = BikeAccident.ParseCsv(this);
	}

	public void draw() {
		background(0);

		pushMatrix();
		translate(width/2, 0);
		// pano.drawThreeFold(width);
		pano.drawTiles(fov, width);
		popMatrix();

		float alpha = 255;
		for(CarAccident accident : carAccidents) {
			stroke(0, alpha);
			fill(0, 100, 0, alpha);
			project(accident.latLng, accident.toString(), 500);
		}

		alpha = 255;
		for(BikeAccident accident : bikeAccidents) {
			stroke(0, alpha);
			fill(0, 100, 0, alpha);
			project(accident.latLng, accident.toString(), 500);
		}

		drawPanoLinks();
		
		pano.drawRoads();

		//Mouse Ref Lin
		stroke(255, 0, 0);
		line(mouseX, 0, mouseX, height);
	}

	public void mousePressed() {
		pov.setHeading(pov.heading()+map(mouseX, 0, width, -fov/2, fov/2));
		pano.jump();
	}

	public void keyPressed() {
		if(key == 'n') {
			pano.jump();
		}
	}

	public void drawPanoLinks() {
		stroke(255);
		PanoLink[] links = pano.getLinks();
		for(int i = 0;i < links.length;i++) {
			int x = pano.headingToPixel(links[i].heading, fov, width);
			line(x, 0, x, height);
		}
	}

	public boolean project(LatLng point, int size) {
		return project(point, null, size);
	}

	public boolean project(LatLng point, String desc, int size) {
		int x = pano.headingToPixel((float)pano.getPosition().getInitialBearing(point), fov, width);
		// line(x, 0, x, height);

		double tanFactor = 65;
		// int horizonHeight = 2*height/5;
		int horizonHeight = height/2;
		double dis = pano.getPosition().getDistance(point);
		if(dis > 50) return false;
		if(dis < 5) return false;
		double vert = Math.atan(dis/tanFactor*Math.PI);
		int y = (int)(height-vert/Math.PI*2*horizonHeight);

		ellipse(x, y, (int)(size/dis), (int)(size*(Math.PI/2-vert)/Math.PI/dis));

		if(desc != null) {
			int padding = 10;
			int xOff = 30;
			int yOff = -130;
			int xSize = 200;
			int ySize = 75;
			int x2 = constrain(x+xOff, 0, width-xSize);
			int y2 = constrain(y+yOff, ySize, height);
			fill(0, 100);
			stroke(0);
			line(x, y, x2, y2);
			rect(x2, y2, xSize, ySize);
			fill(255);
			stroke(255);
			text(desc, x2+padding, y2+padding, xSize-padding, ySize-padding);
		}

		return true;
	}

	public static void main (String [] args) {
		PApplet.main(new String[] { "Ignite" });
	}
}