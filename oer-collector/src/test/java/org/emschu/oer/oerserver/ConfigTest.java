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

import org.junit.Assert;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.emschu.oer.collector.OerCollector;
import org.emschu.oer.core.service.EnvService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.cache.CacheManager;
import org.springframework.context.ApplicationContext;
import org.springframework.scheduling.concurrent.ThreadPoolTaskExecutor;
import org.springframework.test.context.junit4.SpringRunner;

@SpringBootTest(classes = OerCollector.class)
@RunWith(SpringRunner.class)
public class ConfigTest {

    @Autowired
    private ApplicationContext context;

    @Test
    public void testThreadPoolExistence() {
        ThreadPoolTaskExecutor tpte = context.getBean(ThreadPoolTaskExecutor.class);
        Assert.assertNotNull(tpte);
        Assert.assertEquals(tpte.getClass().getName(), ThreadPoolTaskExecutor.class.getName());
    }

    @Test
    public void testCacheManagerExistence() {
        CacheManager cm = context.getBean(CacheManager.class);
        Assert.assertNotNull(cm);
        Assert.assertTrue(cm instanceof CacheManager);
    }

    @Test
    public void testTestModeIsSet() {
        EnvService envService = context.getBean(EnvService.class);
        Assert.assertNotNull(envService);
        Assert.assertEquals(envService.getClass().getName(), EnvService.class.getName());
        Assert.assertTrue(envService.isTestMode());
        Assert.assertFalse(envService.isProdMode());
        Assert.assertFalse(envService.isDevMode());
    }
}
