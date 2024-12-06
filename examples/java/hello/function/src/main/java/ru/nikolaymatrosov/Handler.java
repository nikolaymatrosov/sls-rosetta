package ru.nikolaymatrosov;

import yandex.cloud.sdk.functions.Context;
import yandex.cloud.sdk.functions.YcFunction;

import java.io.PrintStream;
import java.nio.charset.StandardCharsets;

public class Handler implements YcFunction<Integer, String> {
    @Override
    public String handle(Integer i, Context c) {
        System.out.println("stdout: Привет, Мир!");

        PrintStream out = new PrintStream(System.out, true, StandardCharsets.UTF_8);
        out.println("utf8: Привет, Мир!");

        return String.valueOf(i);
    }
}