import processing.core.*;
import panoia.*;

import java.io.*;
import java.util.*;
import java.lang.Math;

public class Ignite extends PApplet {

	Pano pano;
	PanoPov pov;

	float fov = 270;

	ArrayList<CarAccident> carAccidents;
	ArrayList<BikeAccident> bikeAccidents;

	PFont font;
	int roadrand;

	public void setup() {
		size(displayWidth, (int)(displayWidth/(6.5f*fov/360)));
		// size(1024*3, 768);

		background(0);
		smooth();

		font = createFont("Helvetica", 14);
		textFont(font, 14);

		pano = new Pano(this);
		pov = pano.getPov();

		// pano.setPano("3qry8ACTZ8Mw6SQ1UaLNMg");
		// pano.setPosition(new LatLng(45.5110809f, -73.5700496f));
		// pano.setPosition(new LatLng(45.52059937f, -73.58165741f));
		pano.setPano("717wuQJ5lH4xB3Uw5vs4Pw");
		// pano.setPano("FUhF2Lmri2qq6NZErDpn2Q");

		carAccidents = CarAccident.ParseCsv(this);
		bikeAccidents = BikeAccident.ParseCsv(this);

		roadrand = round(random(0, 255));
	}

	public void draw() {
		background(0);

		pushMatrix();
		translate(width/2, 0);
		pano.drawThreeFold(width);
		// pano.drawTiles(fov, width);
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
		drawRoads();

		//Mouse Ref Lin
		// stroke(255, 0, 0);
		// line(mouseX, 0, mouseX, height);
	}

	public void mousePressed() {
		pov.setHeading(pov.heading()+map(mouseX, 0, width, -fov/2, fov/2));
		pano.setPov(pov);

		if (mouseButton == RIGHT) {
			pano.jump();
			roadrand = round(random(0, 255));
		}
	}

	public void keyPressed() {
		if(key == 'n') {
			pano.jump();
			roadrand = round(random(0, 255));
		}
	}

	public void drawPanoLinks() {
		stroke(20);
		strokeWeight(2);
		fill(255, 200);
		PanoLink[] links = pano.getLinks();
		for(int i = 0;i < links.length;i++) {
			int x = pano.headingToPixel(links[i].heading, fov, width);
			// line(x, 0, x, height);
			triangle(x, height-20, x-20, height-7, x+20, height-7);
		}
	}

	public void drawRoads() {
		fill(roadrand, 255-roadrand, 0, 50);
		stroke(255, 50);
		PanoLink[] links = pano.getLinks();
		for(int i = 0;i < links.length;i++) {
			int x = pano.headingToPixel(links[i].heading, fov, width);
			float squish = 6f;
			quad(
				x+width/14/squish, height/2+30, 
				x-width/14/squish, height/2+30, 
				x-width/14, height, 
				x+width/14, height
			);

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