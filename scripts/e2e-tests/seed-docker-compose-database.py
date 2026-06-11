import os
from django.contrib.auth import get_user_model
from django.utils import timezone
from datetime import timedelta
from oauth2_provider.models import AccessToken
from cookbook.models import Recipe, Space, UserSpace, UserPreference
from django_scopes import scope, scopes_disabled
from django.contrib.auth.models import Group

from django.db import connection

User = get_user_model()
with scopes_disabled():
    with connection.cursor() as cursor:
        cursor.execute(f'TRUNCATE TABLE {Recipe._meta.db_table} RESTART IDENTITY CASCADE;')
        cursor.execute(f'TRUNCATE TABLE {UserSpace._meta.db_table} RESTART IDENTITY CASCADE;')
        cursor.execute(f'TRUNCATE TABLE {Space._meta.db_table} RESTART IDENTITY CASCADE;')
        cursor.execute(f'TRUNCATE TABLE {User._meta.db_table} RESTART IDENTITY CASCADE;')

user, _ = User.objects.get_or_create(username='admin', email='admin@example.com')
user.set_password('admin')
user.is_superuser = True
user.is_staff = True
user.save()

from oauth2_provider.models import AccessToken, Application

# Create an Application for the token
app, _ = Application.objects.get_or_create(
    name='Tandoor MCP E2E',
    user=user,
    client_type=Application.CLIENT_CONFIDENTIAL,
    authorization_grant_type=Application.GRANT_PASSWORD,
)

access_token, created = AccessToken.objects.get_or_create(
    user=user,
    token="e2e_test_token",
    defaults={
        "application": app,
        "expires": timezone.now() + timedelta(days=3650),
        "scope": "read write bookmarklet mealplan"
    }
)

from cookbook.helper.permission_helper import create_space_for_user

# Tandoor requires a Space
us = create_space_for_user(user, 'e2e_test_space')
us.active = True
us.save()
space = us.space

# Create a test recipe in the space
with scope(space=space):
    Recipe.objects.get_or_create(name='Tandoori Chicken', created_by=user, space=space, defaults={'description': 'E2E test'})

print("Successfully seeded E2E user, token, space, and recipe.")
