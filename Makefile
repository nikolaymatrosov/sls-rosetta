.PHONY: clean-terraform

clean-terraform:
	find . -path "*/.terraform/providers" -type d -exec rm -rf {} +
