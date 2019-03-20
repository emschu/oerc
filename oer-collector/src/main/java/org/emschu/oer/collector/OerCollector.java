package org.emschu.oer.collector;

/*-
 * #%L
 * oer-server
 * %%
 * Copyright (C) 2019 emschu[aet]mailbox.org
 * %%
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * #L%
 */

import org.springframework.boot.SpringApplication;
import org.springframework.boot.WebApplicationType;
import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.boot.autoconfigure.domain.EntityScan;
import org.springframework.boot.builder.SpringApplicationBuilder;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.PropertySource;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;
import org.springframework.scheduling.annotation.EnableScheduling;

@EnableScheduling
@EnableCaching
@ComponentScan(basePackages = "org.emschu.oer.core")
@ComponentScan(basePackages = "org.emschu.oer.collector")
@EntityScan(basePackages = "org.emschu.oer.core.model")
@PropertySource("classpath:jdbc.properties")
@PropertySource("classpath:oer.properties")
@EnableJpaRepositories(basePackages = "org.emschu.oer.core.model.repository")
@EnableAutoConfiguration
public class OerCollector {

    public static void main(String[] args) {
        // configure this spring boot application
        SpringApplicationBuilder builder = new SpringApplicationBuilder(OerCollector.class)
                .profiles("dev", "test", "prod")
                .parent(OerCollector.class);
        SpringApplication app = builder.build();
        app.setWebApplicationType(WebApplicationType.NONE);
        app.setLogStartupInfo(true);
        app.setHeadless(true);
        app.run(args);
    }
}
