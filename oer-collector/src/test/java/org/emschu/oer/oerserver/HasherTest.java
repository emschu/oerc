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

import junit.framework.TestCase;
import org.junit.Assert;
import org.junit.Test;
import org.emschu.oer.core.util.Hasher;

public class HasherTest extends TestCase {

    @Test
    public void testMd5Hasher() {
        String testString = "test";
        Assert.assertEquals("098F6BCD4621D373CADE4E832627B4F6", Hasher.getHash(testString));
    }

    @Test()
    public void testNullHash() {
        Assert.assertEquals("D41D8CD98F00B204E9800998ECF8427E", Hasher.getHash(""));
        Assert.assertEquals("D41D8CD98F00B204E9800998ECF8427E", Hasher.getHash(null));
    }
}
