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

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;
import org.springframework.context.event.ApplicationEventMulticaster;
import org.springframework.context.event.SimpleApplicationEventMulticaster;
import org.springframework.core.task.SimpleAsyncTaskExecutor;
import org.springframework.scheduling.concurrent.ThreadPoolTaskExecutor;

import java.util.logging.Logger;

@Configuration
@PropertySource(value = "classpath:oer.properties")
public class AsyncEventProcessing {

    private static final Logger LOG = Logger.getLogger(AsyncEventProcessing.class.getName());

    @Value(value = "${oer.collector.core_thread_pool_size}")
    private String coreThreadPoolSize;

    @Value(value = "${oer.collector.max_thread_pool_size}")
    private String maxThreadPoolSize;

    @Bean(name = "applicationEventMulticaster")
    public ApplicationEventMulticaster simpleApplicationEventMulticaster() {
        SimpleApplicationEventMulticaster eventMulticaster
                = new SimpleApplicationEventMulticaster();

        eventMulticaster.setTaskExecutor(new SimpleAsyncTaskExecutor());
        return eventMulticaster;
    }

    @Bean
    public ThreadPoolTaskExecutor threadPoolTaskExecutor() {
        final int maxTpSize;
        final int coreTpSize;
        try {
            maxTpSize = Integer.valueOf(this.maxThreadPoolSize);
            coreTpSize = Integer.valueOf(this.coreThreadPoolSize);
        } catch (NumberFormatException nfe) {
            LOG.warning("invalid thread pool config detected:");
            LOG.warning(this.maxThreadPoolSize);
            LOG.warning(this.coreThreadPoolSize);
            return null;
        }

        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(coreTpSize);
        executor.setWaitForTasksToCompleteOnShutdown(false);
        executor.setThreadNamePrefix("oer_server_thread_pool");
        executor.setMaxPoolSize(maxTpSize);
        executor.setQueueCapacity(Integer.MAX_VALUE);
        executor.initialize();
        return executor;
    }
}

