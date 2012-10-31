import processing.core.*;

import java.io.*;
import java.util.*;

import org.w3c.dom.Document;
import org.w3c.dom.*;

import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.DocumentBuilder;

import java.net.URL;
import java.io.InputStream;

public class Directions {

	public Directions() {

	}

	public Document getXML(LatLng start, LatLng stop) {
		try	{
			DocumentBuilderFactory docBuilderFactory = DocumentBuilderFactory.newInstance();
			DocumentBuilder docBuilder = docBuilderFactory.newDocumentBuilder();

			String urlString = "http://maps.googleapis.com/maps/api/directions/xml?origin="+start.toUrlValue()+"&destination="+stop.toUrlValue()+"&sensor=false&units=metric&mode=driving";
			if(apikey != null && !apikey.equals("")) urlString += "&key="+apikey;

			println(urlString);

			URL url = new URL(urlString);

			InputStream stream = url.openStream();
			Document xml = docBuilder.parse(stream);
			
			return xml;
		} catch (Exception e) {
			System.err.println("Panoia: Error in getXML...");
		}
		
		return null;
	}

	LatLng[] buildRoute(Document xml) {
		try {
			//Parse Location Information and Copyright
			NodeList latTags = xml.getElementsByTagName("lat");
			NodeList lngTags = xml.getElementsByTagName("lng");

			LatLng[] steps = new LatLng[(latTags.getLength()-4)]; //Super hack

			for(int i = 0;i < steps.length;i++) {
				float lat = Float.parseFloat(latTags.item(i).getFirstChild().getNodeValue());
				float lng = Float.parseFloat(lngTags.item(i).getFirstChild().getNodeValue());

				steps[i] = new LatLng(lat, lng);
			}

			println(steps);
			return steps;
		} catch(Exception e) {
			//Probably no such street view, so ignore
			System.err.println("Panoia: Error no such streetview or malformed streetview...");
			return null;
		}

	}

}