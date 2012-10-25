import java.io.*;
import java.util.*;
import processing.core.*;

public class CarAccident {
	
	public Double Latitude;
	public Double Longitude;
	
	private DataVis parent;
	
	public CarAccident(DataVis parent, double lat, double lon) {
		this.parent = parent;
		
		Latitude = lat;
		Longitude = lon;
	}
	
	public void draw(double latitude, double longitude, int pixelWidth, int pixelHeight, double angleWidth, double angleHeight, double headingCenter) {
		
		int scale = 100; // Arbitrarily scaling for testing...
		
		PVector diff = new PVector(scale * (float)(Latitude - latitude), scale * (float)(Longitude - longitude));
		float length = (float)Math.sqrt(Math.pow(diff.x, 2) + Math.pow(diff.y, 2));
		float orientation = (float)(Math.atan2(diff.y, diff.x) * 180 / Math.PI);
		
		float relativeAngle = (float)(orientation - headingCenter);
		while (relativeAngle < 0)
			relativeAngle += 360;
		
		int xPos = (int)(pixelWidth * (relativeAngle / angleWidth)); // This probably doesn't work as expected...
		int yPos = (int)(pixelHeight / 2 + (pixelHeight / (2 * length))); // Approximating horizon to be at half-screen height.
		
		if (yPos > pixelHeight)
			yPos = pixelHeight;
		
		if (length < 1)
			length = 1;
		
		/*
		System.out.println(diff.x + "  " + diff.y);
		System.out.println(xPos);
		System.out.println(yPos);
		*/
		
		System.out.println(length);
		
		parent.stroke(128);
		parent.fill(255);
		parent.ellipseMode(parent.CENTER);
		parent.ellipse(xPos, yPos, 1 * scale / length, 1 * scale / (length * length));
	}
	
	public static ArrayList<CarAccident> ParseCsv(DataVis parent) {
		
		ArrayList<CarAccident> carAccidents = new ArrayList<CarAccident>();
		String[] data = new String[29];
		
		try {
			File file = new File("CarCrashes2007.csv");
			BufferedReader reader = new BufferedReader(new FileReader(file));
			
			String line = reader.readLine(); // Skip the first one with title info
			int col = 0;
			
			while ((line = reader.readLine()) != null) {
				
				StringTokenizer st = new StringTokenizer(line, ";");
				while (st.hasMoreTokens()) {
					data[col] = st.nextToken();
					//data[col] = "4.44444444";
					//System.out.println(st.nextToken());
					col++;
				}
				
				double lat = parseDouble(data[15]);
				double lon = parseDouble(data[16]);
				
				carAccidents.add(new CarAccident(parent, lat, lon));
				col = 0;
			}
		
		} catch (IOException e) {
			e.printStackTrace();
		}
		
		return carAccidents;
	}
	
	private static Double parseDouble(String number) {
		Double dub;
		
		try {
			dub = Double.parseDouble(number);
		} catch (Exception e) {
			dub = 666.666;
		}
		
		return dub;
	}
	
}
