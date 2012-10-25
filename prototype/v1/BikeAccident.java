import panoia.*;
import java.io.*;
import java.util.*;
import processing.core.*;

public class BikeAccident {
	
	public LatLng latLng;
	public String locdesc;
	public int count;
	private PApplet parent;
	
	public BikeAccident(PApplet parent, LatLng latLng, String locdesc) {
		this.parent = parent;
		this.latLng = latLng;

		this.locdesc = locdesc;
		if(locdesc == null) locdesc = "";

		count = 1;
	}
	
	public static ArrayList<BikeAccident> ParseCsv(PApplet parent) {
		ArrayList<BikeAccident> bikeAccidents = new ArrayList<BikeAccident>();
		String[] data = new String[29];
		int errorCount = 0;
		try {
			File file = new File("BikeAccidents.csv");
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
				
				try {
					float lat = Float.parseFloat(data[7]);
					float lng = Float.parseFloat(data[6]);

					LatLng latLng = new LatLng(lat, lng);

					boolean dup = false;

					for(BikeAccident accident : bikeAccidents) {
						if(accident.latLng.equals(latLng)) {
							dup = true;
							accident.count++;
							break;
						}
					}

					if(!dup) {
						String locdesc = data[4].trim().toLowerCase();
						if(data[5].trim().length() != 0) locdesc += "/"+data[5].trim().toLowerCase();

						bikeAccidents.add(new BikeAccident(parent, latLng, locdesc));
					}
				} catch(Exception e) {
					// System.err.println(data[15]+" - "+data[16]);
					// System.err.println("Error parsing CSV for car accidents");
					// System.err.println(e);
					errorCount++;
				}

				col = 0;
			}

		} catch (IOException e) {
			e.printStackTrace();
		}
		
		return bikeAccidents;
	}

	public String toString() {
		String truncDesc = locdesc;
		if(truncDesc.length() > 29) truncDesc = truncDesc.substring(0, 29);
		if(count > 1) return count+" Bike Accidents\n"+locdesc;
		return count+" Bike Accident\n"+locdesc;
	}
}
