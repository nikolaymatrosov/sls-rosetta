package ru.nikolaymatrosov;

import yandex.cloud.sdk.functions.Context;
import yandex.cloud.sdk.functions.YcFunction;

import java.io.IOException;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;

public class Handler implements YcFunction<Integer, String> {
    @Override
    public String handle(Integer i, Context c) {
        System.out.println("stdout: Привет, Мир!");

        PrintStream out = new PrintStream(System.out, true, StandardCharsets.UTF_8);
        out.println("utf8: Привет, Мир!");

        var cl = Thread.currentThread().getContextClassLoader();
        try (var is = cl.getResourceAsStream("ru/nikolaymatrosov/Handler/test.txt")) {
            if (is != null) {
                byte[] content = is.readAllBytes();
                String contentStr = new String(content, StandardCharsets.UTF_8);
                System.out.println(contentStr);
            }
            else {
                System.out.println("ru/nikolaymatrosov/Handler/test.txt");
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }

        return String.valueOf(i);
    }
}