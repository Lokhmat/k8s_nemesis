package ru.hse.edu.asd.vslopatkin.task8.kwic;

import java.util.*;

enum EventType {
    CIRCULATE,
    SORT,
    DISPLAY
}

interface Observer {
    EventType getType();

    void update(Object data);
}


class EventDispatcher {
    private EnumMap<EventType, Observer> observers = new EnumMap<>(EventType.class);

    public void register(Observer observer) {
        observers.put(observer.getType(), observer);
    }

    public void dispatchEvent(EventType eventType, Object data) {
        observers.get(eventType).update(data);
    }
}

class DataLoader {
    private List<String> lines;
    private EventDispatcher eventDispatcher;

    public DataLoader(EventDispatcher eventDispatcher) {
        this.eventDispatcher = eventDispatcher;
    }

    public void loadData(String input) {
        lines = List.of(input.split("\n"));
        notifyObservers(lines);
    }

    private void notifyObservers(List<String> lines) {
        eventDispatcher.dispatchEvent(EventType.CIRCULATE, lines);
    }
}

class Circulator implements Observer {
    private List<String> shiftedLines = new ArrayList<>();
    private EventDispatcher sorterDispatcher;

    public Circulator(EventDispatcher sorterDispatcher) {
        this.sorterDispatcher = sorterDispatcher;
    }

    @Override
    public void update(Object arg) {
        List<String> lines = new ArrayList<>((List<String>) arg);
        for (String line : lines) {
            for (int i = 0; i < line.split(" ").length; i++) {
                shiftedLines.add(shiftLine(line, i));
            }
        }

        sorterDispatcher.dispatchEvent(EventType.SORT, shiftedLines);
    }

    private String shiftLine(String line, int shift) {
        String[] words = line.split(" ");
        return String.join(" ", rotateArray(words, shift));
    }

    private String[] rotateArray(String[] array, int shift) {
        String[] rotated = new String[array.length];
        System.arraycopy(array, shift, rotated, 0, array.length - shift);
        System.arraycopy(array, 0, rotated, array.length - shift, shift);
        return rotated;
    }

    @Override
    public EventType getType() {
        return EventType.CIRCULATE;
    }
}

class Sorter implements Observer {

    private EventDispatcher displayDispatcher;

    public Sorter(EventDispatcher displayDispatcher) {
        this.displayDispatcher = displayDispatcher;
    }

    @Override
    public EventType getType() {
        return EventType.SORT;
    }

    @Override
    public void update(Object arg) {
        List<String> lines = new ArrayList<>((List<String>) arg);
        Collections.sort(lines);
        displayDispatcher.dispatchEvent(EventType.DISPLAY, lines);
    }
}

class Display implements Observer {


    @Override
    public EventType getType() {
        return EventType.DISPLAY;
    }

    @Override
    public void update(Object arg) {
        List<String> lines = new ArrayList<>((List<String>) arg);
        System.out.println("KWIC Results:");
        for (String line : lines) {
            System.out.println(line);
        }
    }
}

public class KWICApp {
    public static void main(String[] args) {
        EventDispatcher eventDispatcher = new EventDispatcher();

        Circulator circulator = new Circulator(eventDispatcher);
        Sorter sorter = new Sorter(eventDispatcher);
        Display display = new Display();

        DataLoader dataLoader = new DataLoader(eventDispatcher);

        eventDispatcher.register(circulator);
        eventDispatcher.register(sorter);
        eventDispatcher.register(display);

        dataLoader.loadData("the quick brown fox\njumps over the lazy dog");
    }
}