package org.emschu.oer.core.service;

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

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.ApplicationContext;
import org.springframework.stereotype.Service;

import java.util.Arrays;

/**
 * this service provides information about app environment
 */
@Service
public class EnvService {

    @Autowired
    private ApplicationContext context;

    public boolean isTestMode() {
        if (Arrays.asList(context.getEnvironment().getActiveProfiles()).contains("test")) {
            return true;
        }
        return false;
    }

    public boolean isDevMode() {
        if (Arrays.asList(context.getEnvironment().getActiveProfiles()).contains("dev")) {
            return true;
        }
        return false;
    }

    public boolean isProdMode() {
        if (Arrays.asList(context.getEnvironment().getActiveProfiles()).contains("prod")) {
            return true;
        }
        return false;
    }

    public String[] getActiveProfiles() {
        return context.getEnvironment().getActiveProfiles();
    }
}
