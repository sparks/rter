import panoia.*;
import java.io.*;
import java.util.*;
import processing.core.*;

public class CarAccident {
	
	public LatLng latLng;
	private PApplet parent;

	String locdesc;

	public int deaths;
	public int seriousInjuries, lightInjuries;

	public int count;
	
	public CarAccident(PApplet parent, LatLng latLng, String locdesc, int deaths, int seriousInjuries, int lightInjuries) {
		this.parent = parent;
		this.latLng = latLng;

		this.locdesc = locdesc;
		if(locdesc == null) locdesc = "";

		this.deaths = deaths;
		this.seriousInjuries = seriousInjuries;
		this.lightInjuries = lightInjuries;

		count = 1;
	}
	
	public static ArrayList<CarAccident> ParseCsv(PApplet parent) {
		ArrayList<CarAccident> carAccidents = new ArrayList<CarAccident>();
		String[] data = new String[29];
		int errorCount = 0;
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
				try {
					float lat = Float.parseFloat(data[15]);
					float lng = Float.parseFloat(data[16]);

					LatLng latLng = new LatLng(lat, lng);

					int deaths = Integer.parseInt(data[17]);
					int seriousInjuries = Integer.parseInt(data[18]);
					int lightInjuries = Integer.parseInt(data[19]);

					boolean dup = false;

					for(CarAccident accident : carAccidents) {
						if(accident.latLng.equals(latLng)) {
							dup = true;
							accident.count++;
							accident.deaths += deaths;
							accident.seriousInjuries += seriousInjuries;
							accident.lightInjuries += lightInjuries;
							break;
						}
					}


					if(!dup) {
						String locdesc = data[8].trim().toLowerCase();
						if(data[11].trim().length() != 0) locdesc += "/"+data[11].trim().toLowerCase();

						carAccidents.add(new CarAccident(parent, latLng, locdesc, deaths, seriousInjuries, lightInjuries));
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
		
		return carAccidents;
	}

	public String toString() {
		String truncDesc = locdesc;
		if(truncDesc.length() > 29) truncDesc = truncDesc.substring(0, 29);
		if(count > 1) return count+" Car Accidents\n"+locdesc+"\nDeaths:"+deaths+" Injuries:"+(lightInjuries+seriousInjuries);
		return count+" Car Accident\n"+locdesc+"\nDeaths:"+deaths+" Injuries:"+(lightInjuries+seriousInjuries);

	}
	
}
