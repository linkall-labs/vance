package com.linkall.sink.mysql.connection;

import com.linkall.sink.mysql.DbConfig;
import com.linkall.sink.mysql.MySQLConfig;
import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;

import java.sql.Connection;
import java.sql.SQLException;
import java.util.Properties;

public class HikariConnectionProvider implements ConnectionProvider {
  private final HikariDataSource dataSource;

  private static final String driverClassName = "com.mysql.cj.jdbc.Driver";

  public static final String JDBC_URL_PATTERN =
      "jdbc:mysql://%s:%s/%s?useInformationSchema=true&nullCatalogMeansCurrent=false&useUnicode=true";
  public static final String CONNECTION_POOL_PREFIX = "connection-pool-";
  public static final int MAX_POOL_SIZE = 5;
  private static final Properties DEFAULT_JDBC_PROPERTIES = initializeDefaultJdbcProperties();

  public HikariConnectionProvider(DbConfig sqlConfig) {
    dataSource = createDataSource(sqlConfig);
  }

  @Override
  public Connection getConnection() throws SQLException {
    return dataSource.getConnection();
  }

  @Override
  public void close() {
    dataSource.close();
  }

  private HikariDataSource createDataSource(DbConfig dbConfig) {
    final HikariConfig config = new HikariConfig();

    String hostName = dbConfig.getHost();

    config.setPoolName(CONNECTION_POOL_PREFIX + hostName + ":" + dbConfig.getPort());
    config.setJdbcUrl(formatJdbcUrl(hostName, dbConfig.getPort(), dbConfig.getDatabase()));
    config.setUsername(dbConfig.getUsername());
    config.setPassword(dbConfig.getPassword());
    config.setMinimumIdle(1);
    config.setMaximumPoolSize(MAX_POOL_SIZE);
    config.setDriverClassName(driverClassName);

    return new HikariDataSource(config);
  }

  private String formatJdbcUrl(String hostName, int port, String database) {
    Properties combinedProperties = new Properties();
    combinedProperties.putAll(DEFAULT_JDBC_PROPERTIES);

    StringBuilder jdbcUrlStringBuilder =
        new StringBuilder(String.format(JDBC_URL_PATTERN, hostName, port, database));

    combinedProperties.forEach(
        (key, value) -> {
          jdbcUrlStringBuilder.append("&").append(key).append("=").append(value);
        });

    return jdbcUrlStringBuilder.toString();
  }

  private static Properties initializeDefaultJdbcProperties() {
    Properties defaultJdbcProperties = new Properties();
    defaultJdbcProperties.setProperty("zeroDateTimeBehavior", "CONVERT_TO_NULL");
    defaultJdbcProperties.setProperty("characterEncoding", "UTF-8");
    defaultJdbcProperties.setProperty("characterSetResults", "UTF-8");
    return defaultJdbcProperties;
  }
}
