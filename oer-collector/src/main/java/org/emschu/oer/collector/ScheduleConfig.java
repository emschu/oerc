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

import org.emschu.oer.collector.reader.ParserException;
import org.emschu.oer.collector.service.UpdaterService;
import org.emschu.oer.core.service.EnvService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.context.annotation.Configuration;
import org.springframework.scheduling.annotation.SchedulingConfigurer;
import org.springframework.scheduling.config.CronTask;
import org.springframework.scheduling.config.ScheduledTaskRegistrar;
import org.springframework.scheduling.support.CronSequenceGenerator;
import org.springframework.stereotype.Component;

import java.util.logging.Logger;

@Configuration
@Component
public class ScheduleConfig implements SchedulingConfigurer {
    private static final Logger LOG = Logger.getLogger(ScheduleConfig.class.getName());

    @Value("${oer.collector.cron_definition}")
    private String cronDefinition;

    @Autowired
    private UpdaterService updaterService;

    @Autowired
    private EnvService envService;

    @Autowired
    private ConfigurableApplicationContext applicationContext;

    @Override
    public void configureTasks(ScheduledTaskRegistrar taskRegistrar) {
        if (cronDefinition == null || cronDefinition.isEmpty() || cronDefinition.equals("null")) {
            LOG.info("Running in single-execution-mode with cron feature disabled");
            try {
                updateData();
            } catch (InterruptedException e) {
                LOG.warning("Update interrupted");
                Thread.currentThread().interrupt();
            } finally {
                LOG.info("Exiting after update");
            }
            if (!envService.isTestMode()) {
                // don't kill context in test env!
                System.exit(0);
            } else {
                LOG.info("No context close, because of test mode");
            }
            return;
        }
        if (!CronSequenceGenerator.isValidExpression(cronDefinition)) {
            LOG.warning(String.format("'%s' is not a valid cron sequence", cronDefinition));
            return;
        }
        LOG.info("Registering a job for expression: " + cronDefinition);

        taskRegistrar.addCronTask(new CronTask(() -> {
                LOG.info("Running in endless cron managed mode");
            try {
                updateData();
            } catch (InterruptedException e) {
                LOG.warning("Update interrupted");
                Thread.currentThread().interrupt();
            }
        }, cronDefinition));
    }

    private void updateData() throws InterruptedException {
        try {
            updaterService.fetchNewTVProgram();
        } catch (ParserException e) {
            LOG.warning("Problem updating data.");
            LOG.throwing(OerCollector.class.getName(), "configure Tasks", e);
        } catch (InterruptedException e) {
            throw e;
        }
    }
}
