package org.emschu.oer.oerserver;

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

import com.fasterxml.classmate.TypeResolver;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.RequestMethod;
import springfox.documentation.builders.PathSelectors;
import springfox.documentation.builders.RequestHandlerSelectors;
import springfox.documentation.builders.ResponseMessageBuilder;
import springfox.documentation.schema.ModelRef;
import springfox.documentation.service.ApiInfo;
import springfox.documentation.service.Contact;
import springfox.documentation.service.ResponseMessage;
import springfox.documentation.spi.DocumentationType;
import springfox.documentation.spring.web.plugins.Docket;
import springfox.documentation.swagger2.annotations.EnableSwagger2;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

/**
 * this configuration bean sets the swagger api configuration to expose
 */
@Configuration
@PropertySource("classpath:application.properties")
@EnableSwagger2
public class SwaggerConfig {

    @Value("${oer.name}")
    private String appName;

    @Value("${oer.latest_api_version}")
    private String latestApiVersion;

    @Autowired
    private TypeResolver typeResolver;

    @Bean
    public Docket apiDocket() {
        ApiInfo apiInfo = new ApiInfo(
                appName,
                "oer-collector - JSON REST server backend for tv program data of public-law channels in Germany, Austria and Switzerland",
                latestApiVersion,
                "",
                new Contact("emschu", "https://github.com/emschu", "look@github.com"),
                "AGPLv3",
                "http://www.gnu.org/licenses/agpl-3.0.de.html",
                new ArrayList<>()
        );

        List<ResponseMessage> defaultResponseMessages = new ArrayList<>();
        defaultResponseMessages.add(new ResponseMessageBuilder()
                .code(500)
                .message("500 message")
                .responseModel(new ModelRef("Error"))
                .build());
        return new Docket(DocumentationType.SWAGGER_2).apiInfo(apiInfo)
                .directModelSubstitute(LocalDate.class, String.class)
                .directModelSubstitute(LocalDateTime.class, String.class)
                .genericModelSubstitutes(ResponseEntity.class)
                .globalResponseMessage(RequestMethod.GET, defaultResponseMessages)
                .enableUrlTemplating(true)
                .additionalModels(typeResolver.resolve(Error.class))
                .select()
                .apis(RequestHandlerSelectors.any())
                .paths(PathSelectors.ant("/v1/**"))
                .build();
    }
}


