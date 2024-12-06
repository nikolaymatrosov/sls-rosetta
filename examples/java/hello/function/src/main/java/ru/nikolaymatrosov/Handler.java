package ru.nikolaymatrosov;

import yandex.cloud.sdk.functions.Context;
import yandex.cloud.sdk.functions.YcFunction;

import java.io.PrintStream;
import java.io.UnsupportedEncodingException;

public class Handler implements YcFunction<Integer, String> {
    @Override
    public String handle(Integer i, Context c) {
        PrintStream out = null;
        try {
            out = new PrintStream(System.out, true, "UTF-8");
        } catch (UnsupportedEncodingException e) {
            throw new RuntimeException(e);
        }
        out.println("Hello, World! Привет, Мир!");
        return String.valueOf(i);
    }
}