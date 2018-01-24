PROVISIONER = provisioner
DRIVER = driver

.PHONY: all
.PHONY: provisioner
.PHONY: driver

all: provisioner driver

provisioner:
	$(MAKE) -C $(PROVISIONER)

driver:
	$(MAKE) -C $(DRIVER)

clean:
	$(MAKE) -C $(PROVISIONER) clean
	$(MAKE) -C $(DRIVER) clean