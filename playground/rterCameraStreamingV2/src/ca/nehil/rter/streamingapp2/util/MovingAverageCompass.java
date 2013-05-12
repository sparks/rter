package ca.nehil.rter.streamingapp2.util;

/**
 * A simple moving average implementation for compass directions.
 * 
 * Gets rid of averaging anomality due to discontinuity from -180 to 180 on the compass.
 *
 * SMA (Simple moving average) sometimes called rolling average, or running average (mean).
 * see: http://en.wikipedia.org/wiki/Moving_average.
 *
 * @author scottkirkwood
 */
public class MovingAverageCompass {
    private float circularBuffer[];
    private float mean;
    private int circularIndex;
    private int count;
    
    private int negativeCount;

    public MovingAverageCompass(int size) {
        circularBuffer = new float[size];
        reset();
    }

    /**
     * Get the current moving average.
     */
    public float getValue() {
        if (negativeCount >= circularBuffer.length) {
        	return -mean;
        } else {
        	return mean;
        }
    }

    /**
     */
    public void pushValue(float x) {
    	if (x < 0) {
    		if(negativeCount < circularBuffer.length) {
    			negativeCount++;
    		}
    	} else {
    		if(negativeCount > 0) {
    			negativeCount--;
    		}
    	}
    	
    	x = Math.abs(x);
    	
        if (count++ == 0) {
            primeBuffer(x);
        }
        float lastValue = circularBuffer[circularIndex];
        mean = mean + (x - lastValue) / circularBuffer.length;
        circularBuffer[circularIndex] = x;
        circularIndex = nextIndex(circularIndex);
    }

    /*
     */
    public void reset() {
        count = 0;
        circularIndex = 0;
        mean = 0;
        negativeCount = 0;
    }

    public long getCount() {
        return count;
    }

    private void primeBuffer(float val) {
        for (int i = 0; i < circularBuffer.length; ++i) {
            circularBuffer[i] = val;
        }
        mean = val;
    }

    private int nextIndex(int curIndex) {
        if (curIndex + 1 >= circularBuffer.length) {
            return 0;
        }
        return curIndex + 1;
    }
}
