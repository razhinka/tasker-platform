package com.engine.config;

import org.flywaydb.core.Flyway;
import javax.sql.DataSource;

public class DatabaseMigration {
    public static void runMigrations(DataSource dataSource) {
        Flyway flyway = Flyway.configure()
                .dataSource(dataSource)
                .baselineOnMigrate(true)
                .load();

        flyway.migrate();
        System.out.println("Миграции БД успешно выполнены!");
    }
}